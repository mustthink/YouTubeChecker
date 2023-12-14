package caller

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/mustthink/YouTubeChecker/config"
)

func New(cfg *config.CallerConfig) (*Caller, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(cfg.YouTubeApiKey))
	if err != nil {
		return nil, err
	}

	call := service.Search.List([]string{"snippet"}).
		Q(cfg.Query).
		MaxResults(cfg.CountResult).
		Order("date").
		Type("video")

	return &Caller{
		service: service,
		call:    call,
		config:  cfg,
	}, nil
}

type Caller struct {
	service *youtube.Service
	call    *youtube.SearchListCall
	config  *config.CallerConfig
}

func (c *Caller) FetchNewVideos() (*youtube.SearchListResponse, error) {
	return c.call.Do()
}
