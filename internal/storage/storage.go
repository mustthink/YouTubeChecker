package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/mustthink/YouTubeChecker/config"
	"github.com/mustthink/YouTubeChecker/internal/notifications"
	"github.com/mustthink/YouTubeChecker/internal/types"
)

var exist = struct{}{}

type (
	VideoStorage struct {
		config *config.SheetConfig
		db     *sql.DB
		sheets *sheets.Service
		videos map[string]struct{}
	}
)

func New(cfg *config.SheetConfig) (*VideoStorage, error) {
	db, err := sql.Open("sqlite3", "internal/storage/storage.db")
	if err != nil {
		return nil, err
	}

	client := getClient(cfg.OauthConfig)
	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("couldn't create sheets service w err: %s", err.Error())
	}

	storage := &VideoStorage{
		config: cfg,
		db:     db,
		sheets: sheetsService,
		videos: make(map[string]struct{}),
	}
	if err := storage.InitDB(); err != nil {
		return nil, fmt.Errorf("couldn't init db w err: %s", err.Error())
	}

	if err := storage.ReadVideos(); err != nil {
		return nil, fmt.Errorf("couldn't read videos w err: %s", err.Error())
	}

	return storage, err
}

func (s *VideoStorage) IsVideoExist(id string) bool {
	_, ok := s.videos[id]
	return ok
}

func (s *VideoStorage) AddNewVideo(v types.Video) error {
	s.videos[v.VideoID] = exist
	if err := s.InsertVideoToDB(v); err != nil {
		return err
	}
	if err := s.writeToSheet(v); err != nil {
		return err
	}
	notifications.Send(v)
	return nil
}
