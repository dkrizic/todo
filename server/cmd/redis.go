package cmd

import (
	"github.com/dkrizic/todo/server/backend"
	"github.com/dkrizic/todo/server/backend/redis"
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

var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Use the in redis backend",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		httpPort, _ := cmd.Flags().GetInt(httpPortFlag)
		grpcPort, _ := cmd.Flags().GetInt(grpcPortFlag)
		redisHost, _ := cmd.Flags().GetString(redisHostFlag)
		redisPort, _ := cmd.Flags().GetInt(redisPortFlag)
		redisUser, _ := cmd.Flags().GetString(redisUserFlag)
		redisPass, _ := cmd.Flags().GetString(redisPassFlag)

		log.WithFields(log.Fields{
			"httpPort":  httpPort,
			"grpcPort":  grpcPort,
			"redisHost": redisHost,
			"redisPort": redisPort,
			"redisUser": redisUser,
		}).Info("Starting redis backend")

		return backend.Backend{
			HttpPort: httpPort,
			GrpcPort: grpcPort,
			Implementation: redis.NewServer(&redis.Config{
				Host: redisHost,
				Port: redisPort,
				User: redisUser,
				Pass: redisPass,
			}),
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

	viper.BindPFlag("TODO_REDIS_HOST", rootCmd.Flags().Lookup(redisHostFlag))
	viper.BindPFlag("TODO_REDIS_PORT", rootCmd.Flags().Lookup(redisPortFlag))
	viper.BindPFlag("TODO_REDIS_USER", rootCmd.Flags().Lookup(redisUserFlag))
	viper.BindPFlag("TODO_REDIS_PASS", rootCmd.Flags().Lookup(redisPassFlag))
}
