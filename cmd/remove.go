package cmd

import (
	"fmt"
	"os"
	"strconv"

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
		idToRemove, _ := strconv.Atoi(args[0])
		cleChan := channelsListElementsByName(channelsChan(cURLs...))
		for cle := range cleChan {
			if cle.ID == idToRemove {
				storage.RemoveChannelURL(cle.Channel.URL)
				fmt.Println(cle.Channel.Name, "removed")
			}
		}
	},
}
