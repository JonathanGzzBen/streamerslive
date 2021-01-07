package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/JonathanGzzBen/streamerslive/pkg/channel"
	"github.com/JonathanGzzBen/streamerslive/pkg/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "displays saved streaming channels",
	Run: func(cmd *cobra.Command, args []string) {
		cURLs, err := storage.ChannelURLs()
		if err != nil {
			fmt.Fprintln(os.Stderr, `No channels stored. Use command "add" to store a channel`)
			return
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "Channel Name", "Stream Title", "Stream URL"})
		channels := make([]channel.Channel, 0)
		cchan := channelsChan(cURLs...)
		for c := range cchan {
			channels = append(channels, c)
		}
		channels = channel.SortByName(channels)
		id := 1
		for _, c := range channels {
			if c.Stream != nil {
				table.Append([]string{strconv.Itoa(id), c.Name, c.Stream.Title, c.Stream.URL})
			} else {
				table.Append([]string{strconv.Itoa(id), c.Name})
			}
			id++
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
