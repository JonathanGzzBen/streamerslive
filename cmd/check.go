package cmd

import (
	"os"
	"sync"
	"time"

	"github.com/JonathanGzzBen/streamerslive/pkg/channel"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var twitchAPICredentials channel.TwitchAPICredentials = channel.TwitchAPICredentials{
	AppAccessToken: "xcqpzgp6lw4araarzpkm0z9gbfgbjo",
	ClientID:       "i9jknyofth9p7zuzkbyxogdglbr9x4",
}

var checkCmd = &cobra.Command{
	Use:   "check [channelURLs...]",
	Short: "check livestreams",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Channel Name", "Stream Title", "Stream URL"})
		cchan := channelsChan(args...)
		for c := range cchan {
			table.Append([]string{c.Name, c.Stream.Title, c.Stream.URL})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func channelsChan(channelURLs ...string) chan channel.Channel {
	cchan := make(chan channel.Channel)
	go func() {
		cclient := channel.NewChannelsClient(channel.TwitchAPICredentials(twitchAPICredentials))

		var wg sync.WaitGroup
		wg.Add(len(channelURLs))
		for _, url := range channelURLs {
			go func(url string) {
				channel, err := cclient.ChannelFromURL(url)
				if err == nil {
					cchan <- *channel
				}
				wg.Done()
			}(url)
			time.Sleep(2 * time.Millisecond)
		}
		wg.Wait()
		close(cchan)
	}()
	return cchan
}
