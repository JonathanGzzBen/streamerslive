package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "streamerslive",
	Short: "StreamersLive is a tool to check streaming channels",
}
