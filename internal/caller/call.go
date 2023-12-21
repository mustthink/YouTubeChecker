package caller

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/mustthink/YouTubeChecker/config"
)

func createCalls(cfg *config.CallerConfig) ([]call, error) {
	calls := make([]call, 0, len(cfg.YouTubeApiKey))
	for _, apiKey := range cfg.YouTubeApiKey {
		service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
		if err != nil {
			return nil, err
		}

		call := newCall(service, cfg.Query, cfg.CountResult)
		calls = append(calls, call)
	}
	return calls, nil
}

type call struct {
	fetchVideo   *youtube.SearchListCall
	fetchChannel *youtube.ChannelsListCall
}

func newCall(service *youtube.Service, query string, count int64) call {
	return call{
		fetchVideo: service.Search.List([]string{"snippet"}).
			Q(query).
			MaxResults(count).
			Order("date").
			Type("video"),
		fetchChannel: service.Channels.List([]string{"statistics"}),
	}
}

func (c call) getNewVideos() (*youtube.SearchListResponse, error) {
	return c.fetchVideo.Do()
}

func (c call) getChannelsInfo(channelID []string) (*youtube.ChannelListResponse, error) {
	return c.fetchChannel.Id(channelID...).Do()
}
