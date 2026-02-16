package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/concerts"
	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/storage"
	"github.com/shakareem/gigoseek/pkg/telegram"
)

func main() {
	storage, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	cfg := config.Get()
	if cfg.TelegramApiToken == "" || cfg.SpotifyClientID == "" || cfg.SpotifyClientSecret == "" {
		log.Fatal("Required private config values not found in configs/private.json")
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	authUpdates := make(chan int64, 100)

	authServer := telegram.NewAuthServer(storage, authUpdates)
	go func() {
		if err := authServer.Run(); err != nil {
			log.Fatal("Failed to start auth server:", err)
		}
	}()

	concertsProvider := &concerts.TimepadConcertProvider{}

	bot := telegram.NewBot(botAPI, storage, concertsProvider, authUpdates)

	bot.Start()
}
