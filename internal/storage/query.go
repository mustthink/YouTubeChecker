package storage

import (
	"fmt"
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

func (s *VideoStorage) InsertVideoToDB(v Video) error {
	_, err := s.db.Exec(insertVideo, v.videoID, v.videoTitle, v.channelID, v.channelTitle,
		v.description, v.publishDate, v.videoURL, v.thumbnailURL)
	return err
}

func (s *VideoStorage) ReadVideos() error {
	rows, err := s.db.Query(readVideos)
	if err != nil {
		return fmt.Errorf("couldn't read videos from DB w err: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var video Video
		err = rows.Scan(&video.videoID)
		if err != nil {
			return fmt.Errorf("couldn't scan videos from db w err: %s", err.Error())
		}

		s.videos[video.videoID] = exist
	}
	return nil
}
