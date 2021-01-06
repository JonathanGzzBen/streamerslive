package storage

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

var channelURLsFilename string

func init() {
	uhd, _ := os.UserHomeDir()
	channelURLsFilename = uhd + "/.streamerslive"
}

// ChannelURLs returns the stored channel URLs
func ChannelURLs() ([]string, error) {
	dat, err := ioutil.ReadFile(channelURLsFilename)
	if err != nil {
		return nil, err
	}
	storedURLs := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(dat)))
	for scanner.Scan() {
		storedURLs = append(storedURLs, scanner.Text())
	}
	return storedURLs, nil
}

// AddChannelURL stores a channel URL
func AddChannelURL(channelURL string) error {
	f, err := os.OpenFile(channelURLsFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(channelURL + "\n"))
	return err
}
