package streams

import (
	"strings"

	"github.com/gocolly/colly"
)

// GetActiveStream returns a channel's active stream
func GetActiveStream(channelURL string) (*Stream, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.youtube.com"),
	)

	channelIsStreaming := false
	// If page has player element, channel is streaming.
	c.OnHTML("#player", func(e *colly.HTMLElement) {
		if e != nil {
			channelIsStreaming = true
		}
	})

	streamTitle := ""
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if channelIsStreaming {
			streamTitle = strings.TrimSpace(strings.TrimSuffix(e.Text, "- YouTube"))
		}
	})

	stream := &Stream{}
	stream.URL = channelURL + "/live"
	c.Visit(stream.URL)
	if !channelIsStreaming {
		return nil, ErrStreamNotActive
	}
	stream.Title = streamTitle
	return stream, nil
}
