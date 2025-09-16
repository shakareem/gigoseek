package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFilePath = "configs/private.json"

type Responses struct {
	Start           string `json:"start"`
	EnterCity       string `json:"enter_city"`
	AuthPrompt      string `json:"auth_prompt"`
	AuthSuccess     string `json:"auth_success"`
	AuthFail        string `json:"auth_fail"`
	Help            string `json:"help"`
	UnknownCommand  string `json:"unknown_command"`
	FavoriteArtists string `json:"favorite_artists"`
}

type Config struct {
	TelegramApiToken    string    `json:"telegram_api_token"`
	SpotifyClientID     string    `json:"spotify_client_id"`
	SpotifyClientSecret string    `json:"spotify_client_secret"`
	TimepadApiToken     string    `json:"timepad_api_token"`
	AuthServerURL       string    `json:"auth_server_url"`
	Responses           Responses `json:"responses"`
}

var cfg *Config

func Get() *Config {
	if cfg == nil {
		var err error
		cfg, err = loadConfig(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
