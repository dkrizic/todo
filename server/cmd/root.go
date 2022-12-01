package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	keyVerbose = "verbose"
)

var verbose = 9

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "The todo backend",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
		postInitCommands(rootCmd.Commands())
	})

	rootCmd.PersistentFlags().CountVarP(&verbose, keyVerbose, "v", "Verbosity (repeat for more verbose output")
	viper.BindEnv(keyVerbose, "TODO_VERBOSE")
	viper.SetDefault(keyVerbose, 4)
	// viper.BindPFlag(keyVerbose, rootCmd.PersistentFlags().Lookup(keyVerbose))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

	fmt.Println("verbose=", verbose)
	switch verbose {
	case 0:
		log.SetLevel(log.ErrorLevel)
		break
	case 1:
		log.SetLevel(log.WarnLevel)
		break
	case 2:
		log.SetLevel(log.InfoLevel)
		break
	case 3:
		log.SetLevel(log.DebugLevel)
		break
	case 4:
		log.SetLevel(log.TraceLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
