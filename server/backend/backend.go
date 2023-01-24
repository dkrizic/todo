package backend

import (
	"fmt"
	repository "github.com/dkrizic/todo/server/backend/repository"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

type Backend struct {
	HttpPort       int
	GrpcPort       int
	HealthPort     int
	MetricsPort    int
	Implementation repository.TodoRepository
	TraceProvider  *tracesdk.TracerProvider
}

var backend Backend

func (backend Backend) Start() (err error) {

	log.WithFields(log.Fields{
		"httpPort":       backend.HttpPort,
		"healthPort":     backend.HealthPort,
		"metricsPort":    backend.MetricsPort,
		"implementation": backend.Implementation,
	}).Info("Starting backend")

	mux := chi.NewMux()
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})
	mux.Handle("/api/v1/todos", otelhttp.NewHandler(http.HandlerFunc(TodosHandler), "todos"))
	mux.Handle("/api/v1/todos/{id}", otelhttp.NewHandler(http.HandlerFunc(TodoHandler), "todos/{id}"))
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
