package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/storage"
	"github.com/shakareem/gigoseek/pkg/telegram"
)

func main() {
	storage := storage.NewInMemoryStorage()

	cfg := config.Get()
	if cfg.TelegramApiToken == "" || cfg.SpotifyClientID == "" || cfg.SpotifyClientSecret == "" {
		log.Fatal("Required private config values not found in configs/private.json")
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Fatal("cannot create bot", err)
	}

	authServer := telegram.NewAuthServer(storage, botAPI)
	go func() {
		if err := authServer.Run(); err != nil {
			log.Fatal("cannot start auth server", err)
		}
	}()

	bot := telegram.NewBot(botAPI, storage)

	err = bot.Start()
	if err != nil {
		log.Fatal("cannot start bot", err)
	}
}
