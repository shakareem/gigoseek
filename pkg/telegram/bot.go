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

const (
	StateIdle storage.ChatState = iota
	StateWaitingForCity
)

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

		chatState, err := b.storage.GetChatState(update.Message.Chat.ID)
		if err != nil {
			chatState = StateIdle
			b.storage.SaveChatState(update.Message.Chat.ID, StateIdle)
		}

		if chatState == StateWaitingForCity {
			err = b.handleCityMessage(update.Message.Chat.ID, update.Message.Text)
			if err != nil {
				log.Printf("Error handling city message: %v", err)
			}
			continue
		}

		err = b.handleMessage(update.Message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	return nil
}

func (b *Bot) handleCityMessage(chatID int64, city string) error {
	// TODO: проверять валидность города (мб через timepad)

	b.storage.SaveCity(chatID, city)
	log.Printf("City for chat %d set successfully", chatID)

	b.storage.SaveChatState(chatID, StateIdle)

	msg := tgbotapi.NewMessage(chatID, "Город успешно установлен!")
	_, err := b.botAPI.Send(msg)

	return err
}
