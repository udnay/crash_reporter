package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crash_reporter",
		Short: "A CLI for copying and exfiltrating core dumps",
	}

	logger *zap.SugaredLogger
)

func init() {
	unsugaredLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger = unsugaredLogger.Sugar()
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
