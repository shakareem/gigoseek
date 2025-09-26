package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFilePath = "configs/config.json"

type Messages struct {
	Start           string `json:"start"`
	AuthPrompt      string `json:"auth_prompt"`
	AuthSuccess     string `json:"auth_success"`
	AuthFail        string `json:"auth_fail"`
	Help            string `json:"help"`
	FavoriteArtists string `json:"favorite_artists"`
	EnterCity       string `json:"enter_city"`
	CitySuccess     string `json:"city_success"`
	NoFavorites     string `json:"no_favorites"`
	NoConcerts      string `json:"no_concerts"`
	WaitForConcerts string `json:"wait_for_concerts"`
}

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"db_name"`
}

type TokensAndSecrets struct {
	TelegramApiToken    string
	SpotifyClientID     string
	SpotifyClientSecret string
	TimepadApiToken     string
}

type TLS struct {
	TLScrtPath string `json:"tls_crt_path"`
	TLSkeyPath string `json:"tls_key_path"`
}

type Timepad struct {
	ApiURL             string `json:"api_url"`
	ConcertsCategoryID string `json:"concerts_category_id"`
}

type Config struct {
	TokensAndSecrets
	AuthServerURL string   `json:"auth_server_url"`
	Messages      Messages `json:"messages"`
	Database      Database `json:"database"`
	TLS           TLS      `json:"tls"`
	Timepad       Timepad  `json:"timepad"`
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

	cfg.TelegramApiToken = os.Getenv("TELEGRAM_API_TOKEN")
	cfg.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	cfg.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")

	return &cfg, nil
}
