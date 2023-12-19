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
	calls := make([]*youtube.SearchListCall, 0, len(cfg.YouTubeApiKey))
	for _, apiKey := range cfg.YouTubeApiKey {
		service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
		if err != nil {
			return nil, err
		}

		call := service.Search.List([]string{"snippet"}).
			Q(cfg.Query).
			MaxResults(cfg.CountResult).
			Order("date").
			Type("video")

		calls = append(calls, call)
	}

	err := os.Mkdir("internal/caller/responses", 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return &Caller{
		calls:     calls,
		config:    cfg,
		responses: 1,
	}, nil
}

type Caller struct {
	calls  []*youtube.SearchListCall
	config *config.CallerConfig

	apiKeyIter int
	responses  int
}

func (c *Caller) FetchNewVideos() (*youtube.SearchListResponse, error) {
	call := c.chooseCall()
	response, err := call.Do()
	if err != nil {
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

	dataToWrite, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil
	}

	if _, err := newResponse.Write(dataToWrite); err != nil {
		return err
	}

	c.responses++
	return newResponse.Close()
}

func (c *Caller) chooseCall() *youtube.SearchListCall {
	defer func() { c.apiKeyIter++ }()
	if c.apiKeyIter == len(c.calls) {
		c.apiKeyIter = 0
	}
	return c.calls[c.apiKeyIter]
}
