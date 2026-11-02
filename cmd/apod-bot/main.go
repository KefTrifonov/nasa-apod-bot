package main

import (
	"log"

	"github.com/KefTrifonov/nasa-apod-bot/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("config loaded, chat_id:", cfg.ChatID)

}
