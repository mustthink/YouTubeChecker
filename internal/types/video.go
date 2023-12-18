package types

import (
	"fmt"

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

func (v Video) Info() string {
	return fmt.Sprintf("Recieved new video:\n\nVideo name: %s\n\nChannel name: %s\n\nLink: \n%s\n",
		v.VideoTitle, v.ChannelTitle, v.VideoURL)
}

func ToVideo(item *youtube.SearchResult, isTracked bool) Video {
	return Video{
		VideoID:      item.Id.VideoId,
		VideoTitle:   item.Snippet.Title,
		ChannelID:    item.Snippet.ChannelId,
		ChannelTitle: item.Snippet.ChannelTitle,
		Description:  item.Snippet.Description,
		PublishDate:  item.Snippet.PublishedAt,
		VideoURL:     fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
		ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		IsTracked:    isTracked,
	}
}
