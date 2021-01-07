package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/JonathanGzzBen/streamerslive/pkg/channel"
	"github.com/JonathanGzzBen/streamerslive/pkg/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [channelId]",
	Short: "removes a channel from list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cURLs, err := storage.ChannelURLs()
		if err != nil {
			fmt.Fprintln(os.Stderr, `No channels stored. Use command "add" to store a channel`)
			return
		}
		channels := make([]channel.Channel, 0)
		cchan := channelsChan(cURLs...)
		for c := range cchan {
			channels = append(channels, c)
		}
		channels = channel.SortByName(channels)
		id := 1
		idToRemove, _ := strconv.Atoi(args[0])
		for _, c := range channels {
			if id == idToRemove {
				storage.RemoveChannelURL(c.URL)
			}
			id++
		}
	},
}
