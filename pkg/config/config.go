package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Responses struct {
	Start             string `json:"start"`
	AlreadyAuthorized string `json:"already_authorized"`
	// TODO others
}

type Config struct {
	TelegramApiToken   string    `json:"telegram_api_token"`
	SpotifyClientID    string    `json:"spotify_client_id"`
	SpotifyClientToken string    `json:"spotify_client_token"`
	AuthServerURL      string    `json:"auth_server_url"`
	Responses          Responses `json:"responses"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
