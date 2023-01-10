package cmd

import (
	"github.com/dkrizic/todo/server/sender"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	httpPortFlag                 = "http-port"
	grpcPortFlag                 = "grpc-port"
	healthPortFlag               = "health-port"
	metricsPortFlag              = "metrics-port"
	notificationsEnabledFlag     = "sender-enabled"
	notificationsPubSubNameFlag  = "sender-pubsub-name"
	notificationsPubSubTopicFlag = "sender-pubsub-topic"
)

var senderClient *sender.Sender

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the actual service",
	Long: `There are multiple backends available for the service like
memory, redis, etc. This command will start the service with the
given backend.`,
	ValidArgs: []string{"memory", "redis"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		enabled := viper.GetBool(notificationsEnabledFlag)
		pubsubName := viper.GetString(notificationsPubSubNameFlag)
		topicName := viper.GetString(notificationsPubSubTopicFlag)
		llog := log.WithFields(log.Fields{
			"enabled":    enabled,
			"pubsubName": pubsubName,
			"topicName":  topicName,
		})
		llog.Info("Creating sender client")
		var err error
		senderClient, err = sender.NewSender(
			enabled,
			pubsubName,
			topicName,
		)
		if err != nil {
			llog.WithError(err).Warn("Unable to create sender client")
			return err
		}
		if senderClient.Enabled {
			senderClient.SendNotification(([]byte)("Hello from Dapr"))
		}
		return nil
	},
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
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
	viper.BindEnv(notificationsEnabledFlag, "TODO_NOTIFICATIONS_ENABLED")
	viper.BindEnv(notificationsPubSubNameFlag, "TODO_NOTIFICATIONS_PUBSUB_NAME")
	viper.BindEnv(notificationsPubSubTopicFlag, "TODO_NOTIFICATIONS_PUBSUB_TOPIC")
}
