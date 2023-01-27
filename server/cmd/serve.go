package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	httpPortFlag                 = "http-port"
	grpcPortFlag                 = "grpc-port"
	healthPortFlag               = "health-port"
	metricsPortFlag              = "metrics-port"
	notificationsEnabledFlag     = "sender-enabled"
	notificationsPubSubNameFlag  = "sender-pubsub-name"
	notificationsPubSubTopicFlag = "sender-pubsub-topic"
	tracingEnabledFlag           = "tracing-enabled"
	tracingEndpointFlag          = "tracing-endpoint"
)

var shutdown func(context.Context) error

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the actual service",
	Long: `There are multiple backends available for the service like
memory, redis, etc. This command will start the service with the
given backend.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tracingEnabled := viper.GetBool(tracingEnabledFlag)
		tracingEndpoint := viper.GetString(tracingEndpointFlag)
		log.WithFields(log.Fields{
			"tracingEnabled":  tracingEnabled,
			"tracingEndpoint": tracingEndpoint,
		}).Info("Tracing configuration")
		var err error
		if tracingEnabled {
			shutdown, err = initProvider(tracingEnabled, tracingEndpoint)
			if err != nil {
				log.WithError(err).Fatal("Failed to initialize tracing provider")
				return err
			}
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Shutting down TracerProvider")
		ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*5)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			log.WithError(err).Fatal("Failed to shutdown TracerProvider")
			return err
		}
		return nil
	},
	ValidArgs: []string{"memory", "redis"},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntP(httpPortFlag, "p", 8080, "The port to listen on for HTTP requests")
	serveCmd.PersistentFlags().IntP(grpcPortFlag, "g", 9090, "The port to listen on for gRPC requests")
	serveCmd.PersistentFlags().IntP(healthPortFlag, "c", 8081, "The port to listen on for health requests")
	serveCmd.PersistentFlags().IntP(metricsPortFlag, "m", 8082, "The port to listen on for metrics requests")
	serveCmd.PersistentFlags().BoolP(notificationsEnabledFlag, "n", false, "Enable notifications")
	serveCmd.PersistentFlags().StringP(notificationsPubSubNameFlag, "", "todo-pubsub", "The name of the pubsub component to use for notifications")
	serveCmd.PersistentFlags().StringP(notificationsPubSubTopicFlag, "", "todo", "The name of the topic to use for notifications")
	serveCmd.PersistentFlags().BoolP(tracingEnabledFlag, "t", false, "Enable tracing")
	serveCmd.PersistentFlags().StringP(tracingEndpointFlag, "", "localhost:4317", "The endpoint to send traces to")
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
	viper.BindEnv(notificationsEnabledFlag, "TODO_NOTIFICATIONS_ENABLED")
	viper.BindEnv(notificationsPubSubNameFlag, "TODO_NOTIFICATIONS_PUBSUB_NAME")
	viper.BindEnv(notificationsPubSubTopicFlag, "TODO_NOTIFICATIONS_PUBSUB_TOPIC")
	viper.BindEnv(tracingEnabledFlag, "TODO_TRACING_ENABLED")
	viper.BindEnv(tracingEndpointFlag, "TODO_TRACING_ENDPOINT")
}

func initProvider(tracingEnabled bool, tracingEndpoint string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("todo"),
		),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to create resource")
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	log.WithField("oltpEndpoint", tracingEndpoint).Info("Connecting to OpenTelemetry Collector")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, tracingEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	log.WithField("oltpEndpoint", tracingEndpoint).Info("Connected to OpenTelemetry Collector")

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.WithError(err).Fatal("Failed to create trace exporter")
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
