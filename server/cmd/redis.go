package cmd

import (
	"github.com/dkrizic/todo/server/backend"
	"github.com/dkrizic/todo/server/memory"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	redisHostFlag = "redis-host"
	redisPortFlag = "redis-port"
	redisUserFlag = "redis-user"
	redisPassFlag = "redis-pass"
)

// memoryCmd represents the memory command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Use the in redis backend",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		httpPort, _ := serveCmd.PersistentFlags().GetInt(httpPortFlag)
		grpcPort, _ := serveCmd.PersistentFlags().GetInt(grpcPortFlag)
		maxEntries, _ := cmd.Flags().GetInt(maxEntriesFlag)
		log.WithFields(log.Fields{
			"httpPort":   httpPort,
			"grpcPort":   grpcPort,
			"maxEntries": maxEntries,
		}).Info("Starting redis backend")
		return backend.Backend{
			HttpPort:       httpPort,
			GrpcPort:       grpcPort,
			Implementation: memory.NewServer(),
		}.Start()
	},
}

func init() {
	serveCmd.AddCommand(redisCmd)

	redisCmd.Flags().String(redisHostFlag, "localhost", "The redis host")
	redisCmd.Flags().Int(redisPortFlag, 6379, "The redis port")
	redisCmd.Flags().String(redisUserFlag, "", "The redis user")
	redisCmd.Flags().String(redisPassFlag, "", "The redis password")

	redisCmd.MarkFlagRequired(redisHostFlag)
	redisCmd.MarkFlagRequired(redisPortFlag)

	viper.BindEnv(redisHostFlag, "TODO_REDIS_HOST")
	viper.BindEnv(redisPortFlag, "TODO_REDIS_PORT")
	viper.BindEnv(redisUserFlag, "TODO_REDIS_USER")
	viper.BindEnv(redisPassFlag, "TODO_REDIS_PASS")
}
