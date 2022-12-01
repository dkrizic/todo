package backend

import (
	"fmt"
	todo "github.com/dkrizic/todo/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
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
	Implementation todo.ToDoServiceServer
}

func (backend Backend) Start() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", backend.GrpcPort))
	if err != nil {
		log.WithField("grpcPort", backend.GrpcPort).WithError(err).Fatal("Failed to listen")
		return err
	}

	s := grpc.NewServer()
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
	log.Fatal(gwServer.ListenAndServe())

	return nil
}