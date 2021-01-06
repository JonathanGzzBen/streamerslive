package cmd

import (
	"log"

	"github.com/JonathanGzzBen/streamerslive/pkg/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [channelURLs...]",
	Short: "adds channels to list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, url := range args {
			err := storage.AddChannelURL(url)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
