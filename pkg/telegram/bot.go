package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/concerts"
	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/storage"
	"golang.org/x/oauth2"
)

type Storage interface {
	SaveAuthState(state string, chatID int64) error
	GetChatIDbyAuthState(state string) (int64, error)
	DeleteAuthState(state string) error

	SaveToken(chatID int64, token oauth2.Token) error
	GetToken(chatID int64) (oauth2.Token, error)
	DeleteToken(chatID int64) error

	SaveCity(chatID int64, city string) error
	GetCity(chatID int64) (string, error)
	DeleteCity(chatID int64) error

	SaveChatState(chatID int64, state storage.ChatState) error
	GetChatState(chatID int64) (storage.ChatState, error)
	DeleteChatState(chatID int64) error
}

type ConcertsProvider interface {
	GetConcerts(artists []string, city string) []concerts.Concert
}

type Bot struct {
	botAPI           *tgbotapi.BotAPI
	storage          Storage
	concertsProvider ConcertsProvider
	authUpdates      <-chan int64
}

func NewBot(botAPI *tgbotapi.BotAPI, storage Storage, concertsProvider ConcertsProvider, authUpdates <-chan int64) *Bot {
	return &Bot{
		botAPI:           botAPI,
		storage:          storage,
		concertsProvider: concertsProvider,
		authUpdates:      authUpdates,
	}
}

const (
	StateIdle storage.ChatState = iota
	StateWaitingForCity
	StateWaitingForAuth
)

func (b *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.botAPI.Send(msg)
	return err
}

func (b *Bot) Start() {
	b.botAPI.Debug = true

	log.Printf("Authorized bot on account %s", b.botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	chatUpdates := b.botAPI.GetUpdatesChan(u)

	for {
		select {
		case update := <-chatUpdates:
			go func() {
				err := b.handleUpdate(update)
				if err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}()
		case chatID := <-b.authUpdates:
			go func() {
				err := b.handleAuthSuccess(chatID)
				if err != nil {
					log.Printf("Error handling auth success: %v", err)
				}
			}()
		}
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	chatState, err := b.storage.GetChatState(update.Message.Chat.ID)
	if err != nil {
		chatState = StateIdle
		b.storage.SaveChatState(update.Message.Chat.ID, StateIdle)
	}

	switch chatState {
	case StateWaitingForAuth:
		return b.handleAuth(update.Message.Chat.ID)
	case StateWaitingForCity:
		return b.handleCityMessage(update.Message.Chat.ID, update.Message.Text)
	}

	return b.handleMessage(update.Message)
}

func (b *Bot) handleAuthSuccess(chatID int64) error {
	log.Printf("Chat %d authed successfully", chatID)
	b.storage.SaveChatState(chatID, StateIdle)
	err := b.sendMessage(chatID, config.Get().Messages.AuthSuccess)
	if err != nil {
		return err
	}

	if !b.isCitySet(chatID) {
		return b.handleSetCity(chatID)
	}

	return nil
}

func (b *Bot) handleCityMessage(chatID int64, city string) error {
	// TODO: проверять валидность города (мб через timepad)

	err := b.storage.SaveCity(chatID, city)
	if err != nil {
		return b.sendMessage(chatID, fmt.Sprintf("%v", err))
	}
	log.Printf("City for chat %d set successfully", chatID)

	b.storage.SaveChatState(chatID, StateIdle)

	return b.sendMessage(chatID, messages.CitySuccess)
}
