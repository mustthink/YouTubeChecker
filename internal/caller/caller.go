package caller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
		service:   service,
		call:      call,
		config:    cfg,
		responses: 1,
	}, nil
}

type Caller struct {
	service *youtube.Service
	call    *youtube.SearchListCall
	config  *config.CallerConfig

	responses int
}

func (c *Caller) FetchNewVideos() (*youtube.SearchListResponse, error) {
	response, err := c.call.Do()
	if err != nil {
		return nil, err
	}

	err = os.Mkdir("internal/caller/responses", 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return response, c.WriteResponse(*response)
}

func (c *Caller) WriteResponse(response youtube.SearchListResponse) error {
	fileName := fmt.Sprintf("internal/caller/responses/response%d.txt", c.responses)
	newResponse, err := os.Create(fileName)
	if err != nil {
		return err
	}

	dataToWrite, err := json.Marshal(response)
	if err != nil {
		return nil
	}

	if _, err := newResponse.Write(dataToWrite); err != nil {
		return err
	}

	c.responses++
	return newResponse.Close()
}
