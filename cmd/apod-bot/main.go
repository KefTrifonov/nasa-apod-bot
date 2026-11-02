package main

import (
	"log"

	"github.com/KefTrifonov/nasa-apod-bot/internal/config"
	"github.com/KefTrifonov/nasa-apod-bot/internal/nasa"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("config loaded, chat_id:", cfg.ChatID)

	apod, err := nasa.Fetch(cfg.NASAKey)
	if err != nil {
		log.Fatalf("Failed to fetch APOD: %v", err)
	}

	log.Println("APOD:", apod.Date, apod.Title, apod.MediaType)
	log.Println("URL:", apod.URL)

}
