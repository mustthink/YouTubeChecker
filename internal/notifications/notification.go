package notifications

import (
	"log"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/mustthink/YouTubeChecker/config"
)

var notificator chan string

func Send(n Notification) {
	notificator <- n.Notification()
}

type (
	Notificator struct {
		tgBot  *tgBot.BotAPI
		config config.NotificationConfig
	}

	Notification interface {
		Notification() string
	}
)

func New(config config.NotificationConfig) (*Notificator, error) {
	bot, err := tgBot.NewBotAPI(config.TelegramBotApiKey)
	if err != nil {
		return nil, err
	}

	return &Notificator{
		tgBot:  bot,
		config: config,
	}, nil
}

func (n *Notificator) Serve() {
	notificator = make(chan string)
	for {
		select {
		case messageText := <-notificator:
			if messageText == "" {
				continue
			}

			message := tgBot.NewMessage(n.config.ReceiverID, messageText)
			if _, err := n.tgBot.Send(message); err != nil {
				log.Println(err)
			}
		}
	}
}
