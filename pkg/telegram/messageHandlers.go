package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/zmb3/spotify/v2"
)

const (
	startCommand      = "start"
	authCommand       = "auth"
	helpCommand       = "help"
	favouritesCommand = "favorites"
)

var responses = config.Get().Responses

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	response := tgbotapi.NewMessage(msg.Chat.ID, "")

	if !msg.IsCommand() {
		response.Text = "пока поддерживаются только команды"
	}

	switch msg.Command() {
	case startCommand:
		response.Text = responses.Start
		_, err := b.botAPI.Send(response)
		if err != nil {
			return err
		}

		if !b.Authorized(msg.Chat.ID) {
			// мб тут надо проверять отдельно token expired
			return b.handleAuth(msg.Chat.ID)
		}
	case authCommand:
		return b.handleAuth(msg.Chat.ID)
	case helpCommand:
		response.Text = responses.Help
	case favouritesCommand:
		return b.handleFavouriteArtists(msg.Chat.ID)
	default:
		response.Text = responses.UnknownCommand
	}

	_, err := b.botAPI.Send(response)
	return err
}

func (b *Bot) handleFavouriteArtists(chatID int64) error {
	if !b.Authorized(chatID) {
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
		response := tgbotapi.NewMessage(chatID, "No favourite artists found")
		_, err = b.botAPI.Send(response)
		return err
	}

	text := responses.FavoriteArtists
	for i, artist := range artistsPage.Artists {
		text += strconv.Itoa(i+1) + ". " + artist.SimpleArtist.Name + "\n"
	}

	response := tgbotapi.NewMessage(chatID, text)
	_, err = b.botAPI.Send(response)
	return err
}

func (b *Bot) Authorized(chatID int64) bool {
	token, err := b.storage.GetToken(chatID)
	return err == nil && token.Valid()
}

func (b *Bot) handleAuth(chatID int64) error {
	state := generateState()

	b.storage.SaveState(state, chatID)

	url := auth.AuthURL(state)
	response := tgbotapi.NewMessage(chatID, responses.AuthPrompt+url)

	_, err := b.botAPI.Send(response)
	return err
}
