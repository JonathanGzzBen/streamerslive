package streamerslive

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	youTubeDomain = "www.youtube.com"
	twitchDomain  = "www.twitch.tv"
)

var (
	// ErrInvalidTwitchAPICredentials is the error thrown
	// when provided invalid Twitch API credentials
	ErrInvalidTwitchAPICredentials = errors.New("invalid Twitch API credentials")
	// ErrInvalidURL is the error thrown when URL is invalid
	ErrInvalidURL = errors.New("invalid URL")
)

// Channel represents a Youtube or Twitch channel
type Channel struct {
	URL  string
	Name string
	// Stream is nil when channel is not streaming
	Stream *Stream
}

// Stream represents a YouTube or Twitch stream
type Stream struct {
	Title string
	URL   string
}

type twitchAPISearchResponse struct {
	Data []struct {
		BroadcasterLanguage string    `json:"broadcaster_language"`
		DisplayName         string    `json:"display_name"`
		GameID              string    `json:"game_id"`
		ID                  string    `json:"id"`
		IsLive              bool      `json:"is_live"`
		TagIds              []string  `json:"tag_ids"`
		ThumbnailURL        string    `json:"thumbnail_url"`
		Title               string    `json:"title"`
		StartedAt           time.Time `json:"started_at"`
	} `json:"data"`
	Pagination struct {
	} `json:"pagination"`
}

// ChannelsClient provides the client instance to get streams info
type ChannelsClient struct {
	twitchAPICredentials TwitchAPICredentials
}

// TwitchAPICredentials are required to get info from Twitch streams.
// Check https://dev.twitch.tv/docs/authentication
type TwitchAPICredentials struct {
	ClientID       string
	AppAccessToken string
}

// NewChannelsClient creates a new ChannelsClient instance with provided credentials
func NewChannelsClient(twitchAPICredentials TwitchAPICredentials) *ChannelsClient {
	return &ChannelsClient{
		twitchAPICredentials: twitchAPICredentials,
	}
}

// ChannelFromURL fetches a channel's information
func (cc *ChannelsClient) ChannelFromURL(url string) (*Channel, error) {
	if strings.HasPrefix(url, "https://"+youTubeDomain) {
		return youtubeChannel(url)
	} else if strings.HasPrefix(url, "https://"+twitchDomain) {
		return twitchChannel(url, cc.twitchAPICredentials)
	}
	return nil, ErrInvalidURL
}

func youtubeChannel(url string) (*Channel, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(youTubeDomain),
	)

	channelName := ""
	c.OnHTML(`link[itemprop="name"]`, func(e *colly.HTMLElement) {
		channelName = e.Attr("content")
	})
	c.Visit(url)

	as, err := activeYouTubeStream(url)
	if err != nil {
		return nil, err
	}
	channel := &Channel{
		URL:  url,
		Name: channelName,
	}
	if as != nil {
		channel.Stream = as
	}
	return channel, nil
}

func twitchChannel(url string, twitchCredentials TwitchAPICredentials) (*Channel, error) {
	channelName := strings.TrimPrefix(url, "https://"+twitchDomain+"/")
	if strings.Contains(channelName, "/") || len(channelName) <= 0 {
		return nil, ErrInvalidURL
	}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/search/channels?query="+channelName, nil)
	if err != nil {
		return nil, ErrInvalidURL
	}
	req.Header.Set("Authorization", "Bearer "+twitchCredentials.AppAccessToken)
	req.Header.Set("Client-Id", twitchCredentials.ClientID)
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == 401 {
		return nil, ErrInvalidTwitchAPICredentials
	}
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	sr := twitchAPISearchResponse{}
	json.Unmarshal(body, &sr)

	if len(sr.Data) <= 0 {
		return nil, ErrInvalidURL
	}

	stream := &Stream{}
	if sr.Data[0].IsLive {
		stream = &Stream{
			URL:   url,
			Title: sr.Data[0].Title,
		}
	} else {
		stream = nil
	}

	return &Channel{
		Name:   sr.Data[0].DisplayName,
		URL:    url,
		Stream: stream,
	}, nil
}

// activeYoutubeStream fetches a channel's active stream,
// returns nil if channel is not streaming
func activeYouTubeStream(channelURL string) (*Stream, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(youTubeDomain),
	)

	channelIsStreaming := false
	// If page has player element, channel is streaming.
	c.OnHTML("#player", func(e *colly.HTMLElement) {
		if e != nil {
			channelIsStreaming = true
		}
	})

	streamTitle := ""
	c.OnHTML(`meta[name="title"]`, func(e *colly.HTMLElement) {
		if channelIsStreaming {
			streamTitle = e.Attr("content")
		}
	})

	stream := &Stream{}
	stream.URL = channelURL + "/live"
	c.Visit(stream.URL)
	if !channelIsStreaming {
		return nil, nil
	}
	stream.Title = streamTitle
	return stream, nil
}

// SortByName returns a slice of channels sorted by name
func SortByName(cs []Channel) []Channel {
	sort.Slice(cs, func(i, j int) bool {
		return strings.Compare(
			strings.ToLower(cs[i].Name),
			strings.ToLower(cs[j].Name)) < 0
	})
	return cs
}
