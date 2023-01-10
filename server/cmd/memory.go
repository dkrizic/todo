package cmd

import (
	"github.com/dkrizic/todo/server/backend"
	"github.com/dkrizic/todo/server/backend/memory"
	"github.com/dkrizic/todo/server/backend/notification"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	maxEntriesFlag = "max-entries"
)

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Use the in memory backend",
	Long: `Be aware that this backend is only in the local process.
It will not work if there are multiple instances running concurrently.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		httpPort, _ := serveCmd.PersistentFlags().GetInt(httpPortFlag)
		grpcPort, _ := serveCmd.PersistentFlags().GetInt(grpcPortFlag)
		healthPort, _ := serveCmd.PersistentFlags().GetInt(healthPortFlag)
		metricsPort, _ := serveCmd.PersistentFlags().GetInt(metricsPortFlag)
		notificationEnabled, _ := serveCmd.PersistentFlags().GetBool(notificationsEnabledFlag)
		maxEntries, _ := cmd.Flags().GetInt(maxEntriesFlag)
		log.WithFields(log.Fields{
			"httpPort":    httpPort,
			"grpcPort":    grpcPort,
			"healthPort":  healthPort,
			"metricsPort": metricsPort,
			"maxEntries":  maxEntries,
		}).Info("Starting memory backend")

		memory := memory.NewServer(maxEntries)

		notification := notification.NewServer(memory, notificationEnabled)

		backend := backend.Backend{
			HttpPort:       httpPort,
			GrpcPort:       grpcPort,
			HealthPort:     healthPort,
			MetricsPort:    metricsPort,
			Implementation: notification,
		}.Start()
		return backend
	},
}

func init() {
	serveCmd.AddCommand(memoryCmd)
	memoryCmd.Flags().IntP(maxEntriesFlag, "", 100, "The maximum number of entries to store in memory")
}
