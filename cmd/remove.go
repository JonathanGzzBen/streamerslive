package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [channelId]",
	Short: "Removes a channel from list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cURLs, err := channelsStorage.ChannelURLs()
		if err != nil {
			fmt.Fprintln(os.Stderr, `No channels stored. Use command "add" to store a channel`)
			return
		}
		idToRemove, _ := strconv.Atoi(args[0])
		cleChan := channelsListElementsByName(channelsChan(cURLs...))
		for cle := range cleChan {
			if cle.ID == idToRemove {
				channelsStorage.RemoveChannelURL(cle.Channel.URL)
				fmt.Println(cle.Channel.Name, "removed")
			}
		}
	},
}
