package backend

import (
	"fmt"
	"github.com/dkrizic/todo/api/todo"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

type Backend struct {
	HttpPort       int
	GrpcPort       int
	HealthPort     int
	MetricsPort    int
	Implementation todo.ToDoServiceServer
	TraceProvider  *tracesdk.TracerProvider
}

func (backend Backend) Start() (err error) {

	log.WithFields(log.Fields{
		"httpPort":       backend.HttpPort,
		"grpcPort":       backend.GrpcPort,
		"healthPort":     backend.HealthPort,
		"metricsPort":    backend.MetricsPort,
		"implementation": backend.Implementation,
	}).Info("Starting backend")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", backend.GrpcPort))
	if err != nil {
		log.WithField("grpcPort", backend.GrpcPort).WithError(err).Fatal("Failed to listen")
		return err
	}

	s := grpc.NewServer()
	grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()))
	grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()))
	log.Info("Tracing enabled")

	todo.RegisterToDoServiceServer(s, backend.Implementation)
	reflection.Register(s)
	log.WithField("grpcPort", backend.GrpcPort).Info("Serving gRPC")
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	log.WithField("grpcPort", backend.GrpcPort).Info("Starting gRPC gateway")
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("127.0.0.1:%d", backend.GrpcPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial", err)
		return err
	}

	gwmux := runtime.NewServeMux()
	err = todo.RegisterToDoServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal("Failed to register gateway", err)
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))))

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", backend.HttpPort),
		Handler: mux,
	}

	log.WithField("httpPort", backend.HttpPort).Info("Serving HTTP and gRPC gateway")
	go func() {
		log.Fatal(gwServer.ListenAndServe())
	}()

	metricsmux := http.NewServeMux()
	metricsmux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", backend.MetricsPort),
		Handler: metricsmux,
	}
	go func() {
		log.Fatal(metricsServer.ListenAndServe())
	}()
	log.WithField("metricsPort", backend.MetricsPort).Info("Serving metrics")

	healthmux := http.NewServeMux()
	healthmux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	healthServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", backend.HealthPort),
		Handler: healthmux,
	}
	log.WithField("healthPort", backend.HealthPort).Info("Serving health")
	log.Fatal(healthServer.ListenAndServe())

	return nil
}
