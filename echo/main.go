package main

import (
	"context"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	todo "github.com/dkrizic/todo/api/todo"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000"
	oltpEndpoint  = "otel-collector.observability:4317"
)

func main() {
	log.Info("Starting app")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize the OTLP exporter and the corresponding trace and metric providers.
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	handler := http.Handler(http.DefaultServeMux)
	handler = otelhttp.NewHandler(handler, "echo")

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.Handle("/test", otelhttp.NewHandler(http.HandlerFunc(TestHandler), "test"))
	r.Handle("/notification", otelhttp.NewHandler(http.HandlerFunc(NotificationHandler), "notification"))
	r.Use(muxlogrus.NewLogger().Middleware)
	http.Handle("/", r)

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

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	log.WithField("oltpEndpoint", oltpEndpoint).Info("Connecting to OpenTelemetry Collector")
	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, oltpEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	log.WithField("oltpEndpoint", oltpEndpoint).Info("Connected to OpenTelemetry Collector")

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := tracesdk.NewBatchSpanProcessor(traceExporter)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("echo").Start(r.Context(), "TestHandler")
	log.WithField("traceparent", span.SpanContext().TraceID()).Info("TestHandler")
	defer span.End()
	dumpHeaders(r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("echo").Start(r.Context(), "NotificationHandler")
	defer span.End()
	dumpHeaders(r)
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
