package streams

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	youTubeDomain = "www.youtube.com"
	twitchDomain  = "www.twitch.tv"
)

// Stream represents a YouTube or Twitch stream
type Stream struct {
	Title       string
	URL         string
	ChannelName string
}

var (
	// ErrStreamNotActive is the error thrown if looking
	// for information of a stream that is not active
	ErrStreamNotActive = errors.New("channel is not streaming")
	// ErrInvalidTwitchAPICredentials is the error thrown
	// when provided invalid Twitch API credentials
	ErrInvalidTwitchAPICredentials = errors.New("invalid Twitch API credentials")
	// ErrInvalidURL is the error thrown when URL is invalid
	ErrInvalidURL = errors.New("channel not found")
)

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

// Client provides the client instance to get streams info
type Client struct {
	twitchAPICredentials TwitchAPICredentials
}

// TwitchAPICredentials are required to get info from Twitch streams
// check https://dev.twitch.tv/docs/authentication
type TwitchAPICredentials struct {
	ClientID       string
	AppAccessToken string
}

// NewStreamsClient creates a new StreamsClient instance with provided credentials
func NewStreamsClient(twitchAPICredentials TwitchAPICredentials) *Client {
	return &Client{
		twitchAPICredentials: twitchAPICredentials,
	}
}

// ActiveStream returns a channel's active stream
func (sc *Client) ActiveStream(channelURL string) (*Stream, error) {
	if strings.HasPrefix(channelURL, "https://"+youTubeDomain) {
		return activeYoutubeStream(channelURL)
	} else if strings.HasPrefix(channelURL, "https://"+twitchDomain) {
		return activeTwitchStream(
			channelURL,
			sc.twitchAPICredentials,
		)
	}
	return nil, errors.New("not valid URL")
}

func activeYoutubeStream(channelURL string) (*Stream, error) {
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
		return nil, ErrStreamNotActive
	}
	stream.Title = streamTitle
	stream.ChannelName, _ = youtubeChannelName(channelURL)
	return stream, nil
}

func youtubeChannelName(channelURL string) (string, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(youTubeDomain),
	)

	channelName := ""
	c.OnHTML(`link[itemprop="name"]`, func(e *colly.HTMLElement) {
		channelName = e.Attr("content")
	})
	c.Visit(channelURL)
	return channelName, nil
}

func activeTwitchStream(channelURL string, twitchCredentials TwitchAPICredentials) (*Stream, error) {
	channelName := strings.TrimLeft(channelURL, "https://"+twitchDomain+"/")
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

	if !sr.Data[0].IsLive {
		return nil, ErrStreamNotActive
	}

	return &Stream{
		Title:       sr.Data[0].Title,
		URL:         channelURL,
		ChannelName: sr.Data[0].DisplayName,
	}, nil
}
