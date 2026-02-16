package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	concertsCommand   = "concerts"
)

// TODO: мб интерфейсы ArtistsProvider и ConcertProvider

var messages = config.Get().Messages

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		return b.sendMessage(msg.Chat.ID, messages.Help)
	}

	switch msg.Command() {
	case startCommand:
		return b.handleStart(msg.Chat.ID)
	case authCommand:
		return b.handleAuth(msg.Chat.ID)
	case helpCommand:
		return b.sendMessage(msg.Chat.ID, messages.Help)
	case favouritesCommand:
		return b.handleFavouriteArtists(msg.Chat.ID)
	case changeCityCommand:
		return b.handleSetCity(msg.Chat.ID)
	case concertsCommand:
		return b.handleConcerts(msg.Chat.ID)
	default:
		return b.sendMessage(msg.Chat.ID, messages.Help)
	}
}

func (b *Bot) handleStart(chatID int64) error {
	err := b.sendMessage(chatID, messages.Start)
	if err != nil {
		return err
	}

	if !b.isAuthorized(chatID) {
		b.storage.SaveChatState(chatID, StateWaitingForAuth)
		return b.handleAuth(chatID)
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
	state := generateState() // мб генерировать не тут, а при появлении нового пользователя
	b.storage.SaveAuthState(state, chatID)

	log.Printf("Generated state %v for chat %v", state, chatID)

	url := auth.AuthURL(state)
	return b.sendMessage(chatID, messages.AuthPrompt+url)
}

func (b *Bot) handleSetCity(chatID int64) error {
	err := b.sendMessage(chatID, messages.EnterCity)
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

	names, err := b.getFavoriteArtistsNames(chatID)
	if err != nil {
		return fmt.Errorf("failed to get favorite artists for chat %d: %w", chatID, err)
	}

	if len(names) == 0 {
		return b.sendMessage(chatID, messages.NoFavorites)
	}

	text := messages.FavoriteArtists
	for i, name := range names {
		text += fmt.Sprintf("%d. %s\n", i+1, name)
	}

	return b.sendMessage(chatID, text)
}

func (b *Bot) handleConcerts(chatID int64) error {
	// тут пока предпологаем, что пользователь аутентифицирован и город установлен
	city, err := b.storage.GetCity(chatID)
	if err != nil {
		return fmt.Errorf("failed to get city for chat %d: %w", chatID, err)
	}

	artists, err := b.getFavoriteArtistsNames(chatID)
	if err != nil {
		return fmt.Errorf("failed to get favorite artists for chat %d: %w", chatID, err)
	}

	if len(artists) == 0 {
		return b.sendMessage(chatID, messages.NoFavorites)
	}

	err = b.sendMessage(chatID, messages.WaitForConcerts)
	if err != nil {
		return err
	}

	concerts := b.concertsProvider.GetConcerts(artists, city)

	if len(concerts) == 0 {
		return b.sendMessage(chatID, messages.NoConcerts)
	}

	var sBuilder strings.Builder
	sBuilder.WriteString(fmt.Sprintf("Найдено %d событий:\n\n", len(concerts)))
	for _, c := range concerts {
		sBuilder.WriteString(fmt.Sprintf("Название: %s\nВремя начала:%s\nСсылка: %s\n\n",
			c.Name,
			c.StartsAt,
			c.URL))
	}

	return b.sendMessage(chatID, sBuilder.String())
}

func (b *Bot) getFavoriteArtistsNames(chatID int64) ([]string, error) {
	log.Printf("Getting top artists for chat ID: %d", chatID)

	token, err := b.storage.GetToken(chatID)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get token: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := spotify.New(auth.Client(ctx, &token))

	artistsPage, err := client.CurrentUsersTopArtists(ctx, spotify.Limit(50))
	if err != nil {
		return []string{}, fmt.Errorf("failed to get top artists: %w", err)
	}

	log.Printf("Retrieved %d artists", len(artistsPage.Artists))

	artistsNames := make([]string, len(artistsPage.Artists))
	for i, artist := range artistsPage.Artists {
		// TODO: мб фильтровать только артистов из России (как?)
		artistsNames[i] = artist.Name
	}

	return artistsNames, nil
}
