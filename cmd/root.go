package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crash_reporter",
		Short: "A CLI for copying and exfiltrating core dumps",
	}

	logger         *zap.SugaredLogger
	configFilePath string
)

func init() {
	unsugaredLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger = unsugaredLogger.Sugar()

	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "File path to a yaml file of arguments for the command")
}

// Execute executes the root command.
func Execute() error {
	logger.Infof("In root command")
	if configFilePath != "" {
		logger.Infof("Reading config file from %s", configFilePath)
		viper.SetConfigFile(configFilePath)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Fatalf("Unable to read config file %s", configFilePath)
			return err
		}
	}
	return rootCmd.Execute()
}
