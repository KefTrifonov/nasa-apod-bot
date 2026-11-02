package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	ChatIDs  []string
	NASAKey  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		NASAKey:  os.Getenv("NASA_API_KEY"),
	}

	// Список получателей через запятую: TELEGRAM_CHAT_IDS=111,222
	// Для совместимости понимаем и старую переменную TELEGRAM_CHAT_ID.
	raw := os.Getenv("TELEGRAM_CHAT_IDS")
	if raw == "" {
		raw = os.Getenv("TELEGRAM_CHAT_ID")
	}
	for _, id := range strings.Split(raw, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			cfg.ChatIDs = append(cfg.ChatIDs, id)
		}
	}

	if cfg.BotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}
	if len(cfg.ChatIDs) == 0 {
		return nil, errors.New("TELEGRAM_CHAT_IDS is not set")
	}
	if cfg.NASAKey == "" {
		cfg.NASAKey = "DEMO_KEY"
	}

	return cfg, nil
}
