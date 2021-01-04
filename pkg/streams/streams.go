package streams

import "errors"

// Stream represents a YouTube or Twitch stream
type Stream struct {
	Title string
	URL   string
}

var (
	// ErrStreamNotActive is the error thrown if looking
	// for information of a stream that is not active
	ErrStreamNotActive = errors.New("channel is not streaming")
)
