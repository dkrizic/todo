package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	httpPortFlag               = "http-port"
	grpcPortFlag               = "grpc-port"
	healthPortFlag             = "health-port"
	metricsPortFlag            = "metrics-port"
	notificationEnabledFlag    = "notification-enabled"
	notificationPubSubNameFlag = "notification-pubsub-name"
	notificationTopicNameFlag  = "notification-topic-name"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the actual service",
	Long: `There are multiple backends available for the service like
memory, redis, etc. This command will start the service with the
given backend.`,
	ValidArgs: []string{"memory", "redis"},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntP(httpPortFlag, "p", 8080, "The port to listen on for HTTP requests")
	serveCmd.PersistentFlags().IntP(grpcPortFlag, "g", 9090, "The port to listen on for gRPC requests")
	serveCmd.PersistentFlags().IntP(healthPortFlag, "c", 8081, "The port to listen on for health requests")
	serveCmd.PersistentFlags().IntP(metricsPortFlag, "m", 8082, "The port to listen on for metrics requests")
	serveCmd.PersistentFlags().BoolP(notificationEnabledFlag, "n", false, "Enable notifications")
	serveCmd.PersistentFlags().StringP(notificationPubSubNameFlag, "", "pubsub", "The name of the pubsub component to use for notifications")
	serveCmd.PersistentFlags().StringP(notificationTopicNameFlag, "", "todo", "The name of the topic to use for notifications")
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
	viper.BindEnv(notificationEnabledFlag, "TODO_NOTIFICATION_ENABLED")
	viper.BindEnv(notificationPubSubNameFlag, "TODO_NOTIFICATION_PUBSUB_NAME")
	viper.BindEnv(notificationTopicNameFlag, "TODO_NOTIFICATION_TOPIC_NAME")
}
