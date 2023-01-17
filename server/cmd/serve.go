package cmd

import (
	"github.com/dkrizic/todo/server/sender"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

const (
	httpPortFlag                 = "http-port"
	grpcPortFlag                 = "grpc-port"
	healthPortFlag               = "health-port"
	metricsPortFlag              = "metrics-port"
	notificationsEnabledFlag     = "sender-enabled"
	notificationsPubSubNameFlag  = "sender-pubsub-name"
	notificationsPubSubTopicFlag = "sender-pubsub-topic"
	enableTracingFlag            = "enable-tracing"
	tracingUrlFlag               = "tracing-url"
)

var senderClient *sender.Sender
var tracer opentracing.Tracer
var closer io.Closer

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the actual service",
	Long: `There are multiple backends available for the service like
memory, redis, etc. This command will start the service with the
given backend.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tracingEnabled := viper.GetBool(enableTracingFlag)
		tracingUrl := viper.GetString(tracingUrlFlag)
		log.WithFields(log.Fields{
			"tracingEnabled": tracingEnabled,
			"tracingUrl":     tracingUrl,
		}).Info("Tracing configuration")
		if tracingEnabled {
			cfg := jaegercfg.Configuration{
				ServiceName: "todo",
				Sampler: &jaegercfg.SamplerConfig{
					Type:  "const",
					Param: 1,
				},
				Reporter: &jaegercfg.ReporterConfig{
					LogSpans:           true,
					LocalAgentHostPort: tracingUrl,
				},
			}

			var err error
			tracer, closer, err = cfg.NewTracer()
			if err != nil {
				log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
			}
			opentracing.SetGlobalTracer(tracer)
		} else {
			log.Info("Tracing disabled")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if closer != nil {
			closer.Close()
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
	serveCmd.PersistentFlags().BoolP(enableTracingFlag, "t", false, "Enable tracing")
	serveCmd.PersistentFlags().StringP(tracingUrlFlag, "", "http://localhost:14268/api/traces", "The url of the tracing server")
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
	viper.BindEnv(notificationsEnabledFlag, "TODO_NOTIFICATIONS_ENABLED")
	viper.BindEnv(notificationsPubSubNameFlag, "TODO_NOTIFICATIONS_PUBSUB_NAME")
	viper.BindEnv(notificationsPubSubTopicFlag, "TODO_NOTIFICATIONS_PUBSUB_TOPIC")
	viper.BindEnv(enableTracingFlag, "TODO_ENABLE_TRACING")
	viper.BindEnv(tracingUrlFlag, "TODO_TRACING_URL")
}
