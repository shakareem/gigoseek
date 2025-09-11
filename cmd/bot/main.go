package main

import (
	"log"

	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/storage"
	"github.com/shakareem/gigoseek/pkg/telegram"
)

func main() {
	cfg := config.Get()

	storage := storage.NewInMemoryStorage()

	authServer := telegram.NewAuthServer(storage)
	go func() {
		if err := authServer.Run(); err != nil {
			log.Fatal("cannot start auth server", err)
		}
	}()

	bot, err := telegram.NewBot(cfg.TelegramApiToken, storage)
	if err != nil {
		log.Fatal("cannot create bot", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatal("cannot start bot", err)
	}
}
