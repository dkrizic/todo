package cmd

import (
	"github.com/dkrizic/todo/server/backend"
	"github.com/dkrizic/todo/server/backend/memory"
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
		maxEntries, _ := cmd.Flags().GetInt(maxEntriesFlag)
		log.WithFields(log.Fields{
			"httpPort":   httpPort,
			"grpcPort":   grpcPort,
			"maxEntries": maxEntries,
		}).Info("Starting memory backend")
		return backend.Backend{
			HttpPort:       httpPort,
			GrpcPort:       grpcPort,
			Implementation: memory.NewServer(),
		}.Start()
	},
}

func init() {
	serveCmd.AddCommand(memoryCmd)
	memoryCmd.Flags().IntP(maxEntriesFlag, "m", 100, "The maximum number of entries to store in memory")
}
