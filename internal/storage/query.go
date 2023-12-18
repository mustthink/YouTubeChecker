package storage

import (
	"fmt"

	"github.com/mustthink/YouTubeChecker/internal/types"
)

const (
	initDB = `
	CREATE TABLE IF NOT EXISTS videos (
		videoID TEXT PRIMARY KEY,
		videoTitle TEXT,
		channelID TEXT,
		channelTitle TEXT,
		description TEXT,
		publishDate TEXT,
		videoURL TEXT,
		thumbnailURL TEXT
	);
	`

	insertVideo = `INSERT INTO videos (videoID, videoTitle, channelID, channelTitle, description, publishDate, videoURL, thumbnailURL)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	readVideos = `SELECT videoID FROM videos`
)

func (s *VideoStorage) InitDB() error {
	_, err := s.db.Exec(initDB)
	return err
}

func (s *VideoStorage) InsertVideoToDB(v types.Video) error {
	_, err := s.db.Exec(insertVideo, v.VideoID, v.VideoTitle, v.ChannelID, v.ChannelTitle,
		v.Description, v.PublishDate, v.VideoURL, v.ThumbnailURL)
	return err
}

func (s *VideoStorage) ReadVideos() error {
	rows, err := s.db.Query(readVideos)
	if err != nil {
		return fmt.Errorf("couldn't read videos from DB w err: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var video types.Video
		err = rows.Scan(&video.VideoID)
		if err != nil {
			return fmt.Errorf("couldn't scan videos from db w err: %s", err.Error())
		}

		s.videos[video.VideoID] = exist
	}
	return nil
}
