package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/config"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	AuthServerURL string
	responses     config.Responses
}

func NewBot(token string, authURL string, responses config.Responses) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot, AuthServerURL: authURL, responses: responses}, err
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

		b.HandleMessage(update.Message)
	}

	return nil
}
