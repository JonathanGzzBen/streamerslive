package storage

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

// ChannelURLsStorage is used to store ChannelURLs
type ChannelURLsStorage struct {
	Filename string
}

// DefaultStorageFilename returns the default location for storing
// channel URLs.
// On UNIX and MacOS, it's $HOME/.streamerslive
// On Windows, it's %USERPROFILE%/.streamerslive
func DefaultStorageFilename() string {
	uhd, _ := os.UserHomeDir()
	return uhd + "/.streamerslive"
}

// NewChannelURLsStorage returns a new ChannelURLsStorage that
// uses filename to store ChannelURLs.
func NewChannelURLsStorage(filename string) *ChannelURLsStorage {
	return &ChannelURLsStorage{
		Filename: filename,
	}
}

// ChannelURLs returns the stored channel URLs
func (cus *ChannelURLsStorage) ChannelURLs() ([]string, error) {
	dat, err := ioutil.ReadFile(cus.Filename)
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

// AddChannelURL add a channel URL to storage
func (cus *ChannelURLsStorage) AddChannelURL(channelURL string) error {
	f, err := os.OpenFile(cus.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(channelURL + "\n"))
	return err
}

// RemoveChannelURL removes a channel URL from storage
func (cus *ChannelURLsStorage) RemoveChannelURL(channelURL string) error {
	dat, err := ioutil.ReadFile(cus.Filename)
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

	os.Remove(cus.Filename)
	for _, url := range urlsToStore {
		cus.AddChannelURL(url)
	}
	return nil
}
