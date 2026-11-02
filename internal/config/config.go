package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	ChatID   string
	NASAKey  string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		NASAKey:  os.Getenv("NASA_API_KEY"),
	}

	if cfg.BotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}
	if cfg.ChatID == "" {
		return nil, errors.New("TELEGRAM_CHAT_ID is not set")
	}
	if cfg.NASAKey == "" {
		cfg.NASAKey = "DEMO_KEY"
	}

	return cfg, nil
}
