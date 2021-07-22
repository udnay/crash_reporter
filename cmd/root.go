package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crash_reporter",
		Short: "A CLI for copying and exfiltrating core dumps",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
