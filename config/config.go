package config

import (
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("couldn't read configuration file w err: %s", err.Error())
	}

	var configuration Config
	if err := json.Unmarshal(data, &configuration); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal configuration w err: %s", err.Error())
	}

	if err := configuration.fullValidation(); err != nil {
		return nil, fmt.Errorf("couldn't validate config w err: %s", err.Error())
	}

	if configuration.TrackedChannelsFilePath == "" {
		configuration.TrackedChannelsFilePath = DefaultTrackedChannels
	}
	if err := configuration.loadTrackedChannels(); err != nil {
		return nil, fmt.Errorf("couldn't load tracked channels w err: %s", err.Error())
	}

	if configuration.SheetConfig.CredentialsPath == "" {
		configuration.SheetConfig.CredentialsPath = DefaultCredentials
	}
	b, err := os.ReadFile(DefaultCredentials)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %s", err.Error())
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %s", err.Error())
	}
	configuration.SheetConfig.OauthConfig = config

	location, err := time.LoadLocation(configuration.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("couldn't load location from time zone w err: %s", err.Error())
	}
	configuration.TimeZoneLocation = location

	return &configuration, nil
}

type (
	CallerConfig struct {
		CountResult   int64    `json:"count_result"`
		YouTubeApiKey []string `json:"api_keys"`
		Query         string   `json:"query"`
	}

	SheetConfig struct {
		OauthConfig     *oauth2.Config `json:"-"`
		CredentialsPath string         `json:"credentials_path"`
		Name            string         `json:"name"`
		SpreadsheetID   string         `json:"id"`
		SheetID         int64          `json:"sheet_id"`
	}

	NotificationConfig struct {
		ReceiverID        int64  `json:"receiver_id"`
		TelegramBotApiKey string `json:"telegram_api_key"`
	}

	Config struct {
		TrackedChannelsFilePath string             `json:"tracked_channels"`
		RequestInterval         time.Duration      `json:"interval_in_seconds"`
		TimeZone                string             `json:"time_zone"`
		CallerConfig            CallerConfig       `json:"caller"`
		SheetConfig             SheetConfig        `json:"sheet"`
		NotificatorConfig       NotificationConfig `json:"notificator"`

		TimeZoneLocation *time.Location `json:"-"`
		trackedChannels  map[string]bool
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
