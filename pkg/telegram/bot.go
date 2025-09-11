package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/storage"
)

type Bot struct {
	botAPI  *tgbotapi.BotAPI
	storage storage.Storage
}

func NewBot(botAPI *tgbotapi.BotAPI, storage storage.Storage) *Bot {
	return &Bot{botAPI: botAPI, storage: storage}
}

func (b *Bot) Start() error {
	b.botAPI.Debug = true

	log.Printf("Authorized on account %s", b.botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		err := b.handleMessage(update.Message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	return nil
}
