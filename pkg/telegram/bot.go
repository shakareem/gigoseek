package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/storage"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	storage storage.Storage
}

func NewBot(token string, storage storage.Storage) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot, storage: storage}, err
}

func (b *Bot) Start() error {
	b.bot.Debug = true

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

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
