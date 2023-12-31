package internal

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mustthink/YouTubeChecker/config"
	"github.com/mustthink/YouTubeChecker/internal/caller"
	"github.com/mustthink/YouTubeChecker/internal/notifications"
	"github.com/mustthink/YouTubeChecker/internal/storage"
	"github.com/mustthink/YouTubeChecker/internal/types"
)

type App struct {
	config       *config.Config
	caller       *caller.Caller
	videoStorage *storage.VideoStorage
	logger       *logrus.Logger
}

func NewApplication(configPath string, isDebug bool) *App {
	logger := logrus.New()
	logger.Info("start initiation application")
	if isDebug {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.Debug("start creating config")
	config, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("couldn't read config w err: %s", err.Error())
	}
	logger.Debug("successfully created config")

	logger.Debug("start creating caller")
	caller, err := caller.New(&config.CallerConfig)
	if err != nil {
		logger.Fatalf("couldn't create caller w err: %s", err.Error())
	}
	logger.Debug("successfully created caller")

	logger.Debug("start creating storage")
	storage, err := storage.New(&config.SheetConfig)
	if err != nil {
		logger.Fatalf("couldn't create storage w err: %s", err.Error())
	}
	logger.Debug("successfully created storage")

	logger.Debug("start creating notificator")
	notificator, err := notifications.New(config.NotificatorConfig)
	if err != nil {
		logger.Fatalf("couldn't create notificator w err: %s", err.Error())
	}
	logger.Debug("successfully created notificator")
	logger.Debug("notificator start serving")
	go notificator.Serve()

	logger.Info("successfully initiated application")
	return &App{
		config:       config,
		caller:       caller,
		videoStorage: storage,
		logger:       logger,
	}
}

func (a *App) Run() {
	a.logger.Info("start application")
	ticker := time.NewTicker(a.config.Interval())
	for ; ; <-ticker.C {
		a.logger.Debug("start fetch new videos")
		response, err := a.caller.FetchNewVideos()
		if err != nil {
			a.logger.Errorf("couldn't call w err: %s", err.Error())
			continue
		}
		a.logger.Debug("successfully fetched new videos")

		a.logger.Debug("start process response")
		a.ProcessResponse(response)
		a.logger.Debug("successfully processed response")
	}
}

func (a *App) ProcessResponse(videos []types.Video) {
	for _, video := range videos {
		a.logger.Debugf("channelID: %s", video.ChannelID)
		a.logger.Debugf("videoID: %s", video.VideoID)

		if !a.videoStorage.IsVideoExist(video.VideoID) {
			video.IsTracked = a.config.IsChannelTracked(video.ChannelID)
			if err := a.videoStorage.AddNewVideo(video); err != nil {
				a.logger.Errorf("couldn't add new video w err: %s", err.Error())
			}
		}
	}

	if err := a.videoStorage.SortSheet(); err != nil {
		a.logger.Errorf("couldn't sort sheet w err: %s", err.Error())
	}
}
