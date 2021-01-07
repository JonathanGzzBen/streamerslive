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

// RemoveChannelURL removes a channel URL from storage
func RemoveChannelURL(channelURL string) error {
	dat, err := ioutil.ReadFile(channelURLsFilename)
	if err != nil {
		return err
	}
	urlsToStore := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(dat)))
	for scanner.Scan() {
		if scanner.Text() != channelURL {
			urlsToStore = append(urlsToStore, scanner.Text())
		}
	}

	os.Remove(channelURLsFilename)
	for _, url := range urlsToStore {
		AddChannelURL(url)
	}
	return nil
}
