package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/JonathanGzzBen/streamerslive/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch [channelId]",
	Short: "Opens streaming channel in browser",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cURLs, err := channelsStorage.ChannelURLs()
		if err != nil {
			fmt.Fprintln(os.Stderr, `No channels stored. Use command "add" to store a channel`)
			return
		}
		idToWatch, _ := strconv.Atoi(args[0])
		cleChan := channelsListElementsByName(channelsChan(cURLs...))
		for cle := range cleChan {
			if cle.ID == idToWatch {
				if cle.Channel.Stream != nil {
					browser.Open(cle.Channel.Stream.URL)
				} else {
					browser.Open(cle.Channel.URL)
				}
			}
		}
	},
}
