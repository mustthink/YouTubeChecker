package types

import (
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	VideoID      string
	VideoTitle   string
	ChannelID    string
	ChannelTitle string
	Description  string
	PublishDate  string
	VideoURL     string
	ThumbnailURL string
	IsTracked    bool
}

func (v Video) ArrayInterface() []interface{} {
	return []interface{}{v.IsTracked, v.VideoID, v.VideoTitle, v.ChannelID, v.ChannelTitle, v.Description, v.PublishDate, v.VideoURL, v.ThumbnailURL}
}

func (v Video) Notification() string {
	return fmt.Sprintf("Publish date: %s\n\nChannel name: %s\n\n%s\n",
		v.PublishDate, v.ChannelTitle, v.VideoURL)
}

func ToVideo(item *youtube.SearchResult, location *time.Location, isTracked bool) (Video, error) {
	utcTime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return Video{}, err
	}
	targetTime := utcTime.In(location)

	return Video{
		VideoID:      item.Id.VideoId,
		VideoTitle:   item.Snippet.Title,
		ChannelID:    item.Snippet.ChannelId,
		ChannelTitle: item.Snippet.ChannelTitle,
		Description:  item.Snippet.Description,
		PublishDate:  targetTime.Format("02.01.2006 15:04"),
		VideoURL:     fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
		ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		IsTracked:    isTracked,
	}, nil
}
