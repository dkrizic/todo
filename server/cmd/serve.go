package cmd

import (
	"context"
	"github.com/dkrizic/todo/server/sender"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
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
	tracingAgentHostFlag         = "tracing-agent-host"
	tracingAgentPortFlag         = "tracing-agent-port"
)

var senderClient *sender.Sender
var traceProvider *tracesdk.TracerProvider

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the actual service",
	Long: `There are multiple backends available for the service like
memory, redis, etc. This command will start the service with the
given backend.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tracingEnabled := viper.GetBool(tracingEnabledFlag)
		tracingAgentHost := viper.GetString(tracingAgentHostFlag)
		tracingAgentPort := viper.GetString(tracingAgentPortFlag)
		log.WithFields(log.Fields{
			"tracingEnabled":   tracingEnabled,
			"tracingAgentHost": tracingAgentHost,
			"tracingAgentPort": tracingAgentPort,
		}).Info("Tracing configuration")
		if tracingEnabled {
			exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(tracingAgentHost), jaeger.WithAgentPort(tracingAgentPort)))
			if err != nil {
				log.WithError(err).Fatal("Failed to create Jaeger exporter")
				return err
			}
			traceProvider = tracesdk.NewTracerProvider(
				// Always be sure to batch in production.
				tracesdk.WithBatcher(exp),
				// Record information about this application in a Resource.
				tracesdk.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("todo"),
					attribute.String("environment", "test"),
				)),
			)
			otel.SetTracerProvider(traceProvider)
			log.WithFields(log.Fields{
				"tracingAgentHost": tracingAgentHost,
				"tracingAgentPort": tracingAgentPort,
			}).Info("Tracing established")
		} else {
			log.Info("Tracing is disabled")
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*5)
		defer cancel()
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.WithError(err).Fatal("Failed to shutdown tracer provider")
		}
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
	serveCmd.PersistentFlags().StringP(tracingAgentHostFlag, "", "localhost", "The host of the tracing agent")
	serveCmd.PersistentFlags().StringP(tracingAgentPortFlag, "", "6831", "The port of the tracing agent")
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
	viper.BindEnv(notificationsEnabledFlag, "TODO_NOTIFICATIONS_ENABLED")
	viper.BindEnv(notificationsPubSubNameFlag, "TODO_NOTIFICATIONS_PUBSUB_NAME")
	viper.BindEnv(notificationsPubSubTopicFlag, "TODO_NOTIFICATIONS_PUBSUB_TOPIC")
	viper.BindEnv(tracingEnabledFlag, "TODO_TRACING_ENABLED")
	viper.BindEnv(tracingAgentHostFlag, "TODO_TRACING_AGENT_HOST")
	viper.BindEnv(tracingAgentPortFlag, "TODO_TRACING_AGENT_PORT")
}
