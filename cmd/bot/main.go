package main

import (
	"log"

	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/telegram"
)

func main() {
	cfg, err := config.LoadConfig("configs/private.json")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	bot, err := telegram.NewBot(cfg.TelegramApiToken, cfg.AuthServerURL, cfg.Responses)
	if err != nil {
		log.Fatal("cannot create bot", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatal("cannot start bot", err)
	}
}
