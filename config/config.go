package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	DefaultConfig          = "config/default.json"
	DefaultTrackedChannels = "config/tracked_channels.json"
	DefaultCredentials     = "config/credentials.json"
)

func New(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configuration Config
	if err := json.Unmarshal(data, &configuration); err != nil {
		return nil, err
	}

	if configuration.TrackedChannelsFilePath == "" {
		configuration.TrackedChannelsFilePath = DefaultTrackedChannels
	}
	if err := configuration.loadTrackedChannels(); err != nil {
		return nil, err
	}

	if configuration.SheetConfig.CredentialsPath == "" {
		configuration.SheetConfig.CredentialsPath = DefaultCredentials
	}
	b, err := os.ReadFile(DefaultCredentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	configuration.SheetConfig.OauthConfig = config

	return &configuration, nil
}

type (
	CallerConfig struct {
		CountResult   int64  `json:"count_result"`
		YouTubeApiKey string `json:"api_key"`
		Query         string `json:"query"`
	}

	SheetConfig struct {
		OauthConfig     *oauth2.Config `json:"-"`
		CredentialsPath string         `json:"credentials_path"`
		Name            string         `json:"name,omitempty"`
		SpreadsheetID   string         `json:"id,omitempty"`
	}

	Config struct {
		TrackedChannelsFilePath string        `json:"tracked_channels"`
		RequestInterval         time.Duration `json:"interval_in_seconds"`
		CallerConfig            CallerConfig  `json:"caller"`
		SheetConfig             SheetConfig   `json:"sheet"`

		trackedChannels map[string]bool
	}
)

func (c *Config) loadTrackedChannels() error {
	file, err := os.ReadFile(c.TrackedChannelsFilePath)
	if err != nil {
		return err
	}

	var rawChannels []string
	if err := json.Unmarshal(file, &rawChannels); err != nil {
		return err
	}

	channels := make(map[string]bool)
	for _, name := range rawChannels {
		channels[name] = true
	}
	c.trackedChannels = channels

	return nil
}

func (c *Config) Interval() time.Duration {
	return c.RequestInterval * time.Second
}

func (c *Config) IsChannelTracked(id string) bool {
	return c.trackedChannels[id]
}
