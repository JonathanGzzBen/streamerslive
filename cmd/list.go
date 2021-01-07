package cmd

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/JonathanGzzBen/streamerslive"
	"github.com/JonathanGzzBen/streamerslive/pkg/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ChannelListElement is a Channel displayed by list command
type ChannelListElement struct {
	ID      int
	Channel streamerslive.Channel
}

var twitchAPICredentials streamerslive.TwitchAPICredentials = streamerslive.TwitchAPICredentials{
	AppAccessToken: "xcqpzgp6lw4araarzpkm0z9gbfgbjo",
	ClientID:       "i9jknyofth9p7zuzkbyxogdglbr9x4",
}

var channelsStorage = storage.NewChannelURLsStorage(storage.DefaultStorageFilename())

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays saved streaming channels",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cURLs, err := channelsStorage.ChannelURLs()
		if err != nil {
			fmt.Fprintln(os.Stderr, `No channels stored. Use command "add" to store a channel`)
			return
		}
		printChannelsList(channelsListElementsByName(channelsChan(cURLs...)))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// channelsChan retrieves channels data
func channelsChan(channelURLs ...string) chan streamerslive.Channel {
	cchan := make(chan streamerslive.Channel)
	go func() {
		cclient := streamerslive.NewChannelsClient(streamerslive.TwitchAPICredentials(twitchAPICredentials))
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

// channelsListElementsByName sorts channels by name and returns channel list elements
func channelsListElementsByName(cChan <-chan streamerslive.Channel) <-chan ChannelListElement {
	cleChan := make(chan ChannelListElement)
	go func() {
		channels := make([]streamerslive.Channel, 0)
		for c := range cChan {
			channels = append(channels, c)
		}
		channels = streamerslive.SortByName(channels)
		id := int64(1)
		for _, c := range channels {
			cleChan <- ChannelListElement{
				ID:      int(atomic.LoadInt64(&id)),
				Channel: c,
			}
			atomic.AddInt64(&id, 1)
		}
		close(cleChan)
	}()
	return cleChan
}

func printChannelsList(cleChan <-chan ChannelListElement) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Channel", "Stream Title"})
	for cle := range cleChan {
		if cle.Channel.Stream != nil {
			table.Append([]string{
				strconv.Itoa(cle.ID),
				cle.Channel.Name,
				cle.Channel.Stream.Title,
			})
		} else {
			table.Append([]string{
				strconv.Itoa(cle.ID),
				cle.Channel.Name,
				"",
			})
		}
	}
	table.Render()
}
