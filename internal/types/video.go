package types

import (
	"fmt"
)

type Video struct {
	VideoID         string
	VideoTitle      string
	ChannelID       string
	ChannelTitle    string
	Description     string
	PublishDate     string
	FetchDate       string
	VideoURL        string
	ThumbnailURL    string
	IsTracked       bool
	SubscriberCount uint64
}

func (v Video) ArrayInterface() []interface{} {
	return []interface{}{v.IsTracked, v.VideoID, v.VideoTitle, v.Description, v.ChannelID, v.ChannelTitle, v.SubscriberCount, v.PublishDate, v.FetchDate, v.VideoURL, v.ThumbnailURL}
}

func (v Video) Notification() string {
	return fmt.Sprintf("Publish date: %s\n\nChannel name: %s\nSubsribers count: %d\n%s\n",
		v.PublishDate, v.ChannelTitle, v.SubscriberCount, v.VideoURL)
}
