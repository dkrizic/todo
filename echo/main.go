package main

import (
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	todo "github.com/dkrizic/todo/api/todo"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"net/http"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000"
)

func main() {
	log.Info("Starting app")

	handler := http.Handler(http.DefaultServeMux)
	handler = otelhttp.NewHandler(handler, "echo")

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.Handle("/test", otelhttp.NewHandler(http.HandlerFunc(TestHandler), "test"))
	r.Handle("/notification", otelhttp.NewHandler(http.HandlerFunc(NotificationHandler), "notification"))
	r.Use(muxlogrus.NewLogger().Middleware)
	http.Handle("/", r)

	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost("localhost"), jaeger.WithAgentPort("6831")))
	if err != nil {
		log.Fatal(err)
	}
	traceProvider := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("echo"),
			attribute.String("environment", "test"),
		)),
	)
	otel.SetTracerProvider(traceProvider)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		log.Info("Shutting down tracer")
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("listenAddress", listenAddress).Info("Starting listener")
	log.Fatal(srv.ListenAndServe())
	log.Info("Stopping listener")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("echo").Start(r.Context(), "TestHandler")
	defer span.End()
	dumpHeaders(r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("echo").Start(r.Context(), "NotificationHandler")
	defer span.End()
	event, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		log.Print("failed to parse CloudEvent from request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	log.WithFields(log.Fields{
		"time":            event.Time(),
		"source":          event.Source(),
		"type":            event.Type(),
		"subject":         event.Subject(),
		"id":              event.ID(),
		"specversion":     event.SpecVersion(),
		"datacontenttype": event.DataContentType(),
		"dataschema":      event.DataSchema(),
		"data":            string(event.Data()),
	}).Info("Received event")

	var change todo.Change
	json.Unmarshal(event.Data(), &change)
	if change.Before != nil {
		log.WithFields(log.Fields{
			"beforeId":          change.Before.Id,
			"beforeTitle":       change.Before.Title,
			"beforeDescription": change.Before.Description,
		}).Info("Before")
	}
	if change.After != nil {
		log.WithFields(log.Fields{
			"afterId":          change.After.Id,
			"afterTitle":       change.After.Title,
			"afterDescription": change.After.Description,
		}).Info("After")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// methods that dumps all header values from a given request to log
func dumpHeaders(r *http.Request) {
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.WithField("header", name).Info(value)
		}
	}
}
