package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/JonathanGzzBen/streamerslive/pkg/streams"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var twitchAPICredentials streams.TwitchAPICredentials = streams.TwitchAPICredentials{
	AppAccessToken: "xcqpzgp6lw4araarzpkm0z9gbfgbjo",
	ClientID:       "i9jknyofth9p7zuzkbyxogdglbr9x4",
}

var rootCmd = &cobra.Command{
	Use:   "streamerslive [channelURLs...]",
	Short: "StreamersLive is a tool to check livestreams",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Channel Name", "Stream Title", "Stream URL"})
		asc := activeStreamsChan(args...)
		for as := range asc {
			table.Append([]string{as.ChannelName, as.Title, as.URL})
		}
		table.Render()
	},
}

func activeStreamsChan(channelsURLs ...string) chan streams.Stream {
	asc := make(chan streams.Stream)
	go func() {
		sc := streams.NewStreamsClient(twitchAPICredentials)
		var wg sync.WaitGroup
		wg.Add(len(channelsURLs))
		for _, cu := range channelsURLs {
			go func(cu string) {
				stream, err := sc.ActiveStream(cu)
				if err == nil {
					asc <- *stream
				}
				wg.Done()
			}(cu)
			time.Sleep(2 * time.Millisecond)
		}
		wg.Wait()
		close(asc)
	}()
	return asc
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
