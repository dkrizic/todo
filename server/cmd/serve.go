package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	httpPortFlag    = "http-port"
	grpcPortFlag    = "grpc-port"
	healthPortFlag  = "health-port"
	metricsPortFlag = "metrics-port"
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
	viper.BindEnv(httpPortFlag, "TODO_HTTP_PORT")
	viper.BindEnv(grpcPortFlag, "TODO_GRPC_PORT")
	viper.BindEnv(healthPortFlag, "TODO_HEALTH_PORT")
	viper.BindEnv(metricsPortFlag, "TODO_METRICS_PORT")
}
