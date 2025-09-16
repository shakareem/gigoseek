package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

const (
	startCommand      = "start"
	authCommand       = "auth"
	helpCommand       = "help"
	favouritesCommand = "favorites"
	changeCityCommand = "changecity"
)

var messages = config.Get().Messages

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		return b.handleHelpCommand(msg.Chat.ID)
	}

	switch msg.Command() {
	case startCommand:
		return b.handleStart(msg.Chat.ID)
	case authCommand:
		return b.handleAuth(msg.Chat.ID)
	case helpCommand:
		return b.handleHelpCommand(msg.Chat.ID)
	case favouritesCommand:
		return b.handleFavouriteArtists(msg.Chat.ID)
	case changeCityCommand:
		return b.handleSetCity(msg.Chat.ID)
	default:
		return b.handleHelpCommand(msg.Chat.ID)
	}
}

func (b *Bot) handleHelpCommand(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, messages.Help)
	_, err := b.botAPI.Send(msg)
	return err
}

func (b *Bot) handleStart(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, messages.Start)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		return err
	}

	if !b.isAuthorized(chatID) {
		err = b.handleAuth(chatID)
		if err != nil {
			return err
		}
		// TODO: мб ждать пока авторизуется
	}

	if !b.isCitySet(chatID) {
		return b.handleSetCity(chatID)
	}

	return nil
}

func (b *Bot) isCitySet(chatID int64) bool {
	_, err := b.storage.GetCity(chatID)
	return err == nil
}

func (b *Bot) isAuthorized(chatID int64) bool {
	token, err := b.storage.GetToken(chatID)
	if err != nil {
		log.Printf("Failed to get token for chat %d: %v", chatID, err)
		return false
	}

	if !token.Valid() {
		err = b.refreshExpiredToken(chatID, &token)
		if err != nil {
			log.Println(err)
			return false
		}
		return b.isAuthorized(chatID)
	}

	return true
}

func (b *Bot) refreshExpiredToken(chatID int64, token *oauth2.Token) error {
	newToken, err := auth.RefreshToken(context.Background(), token)
	if err != nil {
		return fmt.Errorf("failed to refresh access token for chat %d: %w", chatID, err)
	}
	b.storage.DeleteToken(chatID)
	err = b.storage.SaveToken(chatID, *newToken)
	if err != nil {
		return fmt.Errorf("failed to save refreshed token: %w", err)
	}

	log.Printf("Token for chat %d refreshed successfully", chatID)
	return nil
}

func (b *Bot) handleAuth(chatID int64) error {
	state := generateState()

	b.storage.SaveState(state, chatID)

	url := auth.AuthURL(state)
	msg := tgbotapi.NewMessage(chatID, messages.AuthPrompt+url)

	_, err := b.botAPI.Send(msg)
	return err
}

func (b *Bot) handleSetCity(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, messages.EnterCity)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		return err
	}

	b.storage.SaveChatState(chatID, StateWaitingForCity)
	log.Printf("Set chat %d state to waiting for city", chatID)

	return nil
}

func (b *Bot) handleFavouriteArtists(chatID int64) error {
	if !b.isAuthorized(chatID) {
		return b.handleAuth(chatID)
	}

	log.Printf("Getting favorites for chat ID: %d", chatID)

	token, err := b.storage.GetToken(chatID)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := spotify.New(auth.Client(ctx, &token))

	artistsPage, err := client.CurrentUsersTopArtists(
		ctx,
		spotify.Limit(5),
		spotify.Locale("ru_RU"),
	)
	if err != nil {
		return fmt.Errorf("failed to get top artists: %w", err)
	}

	log.Printf("Retrieved %d artists", len(artistsPage.Artists))

	if len(artistsPage.Artists) == 0 {
		msg := tgbotapi.NewMessage(chatID, "No favourite artists found")
		_, err = b.botAPI.Send(msg)
		return err
	}

	text := messages.FavoriteArtists
	for i, artist := range artistsPage.Artists {
		text += fmt.Sprintf("%d. %s\n", i+1, artist.SimpleArtist.Name)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.botAPI.Send(msg)
	return err
}
