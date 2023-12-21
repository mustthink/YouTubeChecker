package caller

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/youtube/v3"

	"github.com/mustthink/YouTubeChecker/config"
	"github.com/mustthink/YouTubeChecker/internal/types"
)

func New(cfg *config.CallerConfig) (*Caller, error) {
	calls, err := createCalls(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't create fetch video calls w err: %s", err.Error())
	}

	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("couldn't load location from time zone w err: %s", err.Error())
	}

	err = os.Mkdir("internal/caller/responses", 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("couldn't create response folder w err: %s", err.Error())
	}

	return &Caller{
		calls:      calls,
		tzLocation: location,
		config:     cfg,
		responses:  1,
	}, nil
}

type Caller struct {
	calls      []call
	config     *config.CallerConfig
	tzLocation *time.Location

	apiKeyIter int
	responses  int
}

func (c *Caller) FetchNewVideos() ([]types.Video, error) {
	call := c.getCall()
	response, err := call.getNewVideos()
	if err != nil {
		return nil, fmt.Errorf("couldn't get new videos w err: %s", err.Error())
	}
	if err := c.writeResponse(*response); err != nil {
		return nil, fmt.Errorf("couldn't write response w err: %s", err.Error())
	}

	channelsIDs := channelIDs(response)
	channelsInfo, err := call.getChannelsInfo(channelsIDs)
	if err != nil {
		return nil, fmt.Errorf("couldn't get channels info w err: %s", err.Error())
	}
	subsCount := subscriberCount(channelsInfo)

	return convertResponseToVideos(response, c.tzLocation, subsCount)
}

func (c *Caller) writeResponse(response youtube.SearchListResponse) error {
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

func (c *Caller) getCall() call {
	defer func() { c.apiKeyIter++ }()
	if c.apiKeyIter == len(c.calls) {
		c.apiKeyIter = 0
	}
	return c.calls[c.apiKeyIter]
}

func channelIDs(list *youtube.SearchListResponse) []string {
	channelIDs := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		channelIDs = append(channelIDs, item.Snippet.ChannelId)
	}
	return channelIDs
}

func subscriberCount(list *youtube.ChannelListResponse) map[string]uint64 {
	count := make(map[string]uint64)
	for _, item := range list.Items {
		count[item.Id] = item.Statistics.SubscriberCount
	}
	return count
}

func convertResponseToVideos(response *youtube.SearchListResponse, location *time.Location, countList map[string]uint64) ([]types.Video, error) {
	videos := make([]types.Video, 0, len(response.Items))
	for _, item := range response.Items {
		utcTime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse time w err: %s", err.Error())
		}
		targetTime := utcTime.In(location)
		nowTime := time.Now().In(location)

		video := types.Video{
			VideoID:         item.Id.VideoId,
			VideoTitle:      item.Snippet.Title,
			ChannelID:       item.Snippet.ChannelId,
			ChannelTitle:    item.Snippet.ChannelTitle,
			Description:     item.Snippet.Description,
			PublishDate:     targetTime.Format("02.01.2006 15:04"),
			FetchDate:       nowTime.Format("02.01.2006 15:04"),
			VideoURL:        fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
			ThumbnailURL:    item.Snippet.Thumbnails.High.Url,
			SubscriberCount: countList[item.Snippet.ChannelId],
		}
		videos = append(videos, video)
	}
	return videos, nil
}
