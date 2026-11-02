package main

import (
	"log"
	"os"
	"strconv"

	"github.com/KefTrifonov/nasa-apod-bot/internal/config"
	"github.com/KefTrifonov/nasa-apod-bot/internal/nasa"
	"github.com/KefTrifonov/nasa-apod-bot/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("config loaded, recipients:", cfg.ChatIDs)

	// Необязательный аргумент: номер записи в ленте (0 — сегодня, 1 — вчера...)
	index := 0
	if len(os.Args) > 1 {
		index, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Bad argument %q: expected a number (0 = today)", os.Args[1])
		}
	}

	apod, err := nasa.Fetch(index)
	if err != nil {
		log.Fatalf("Failed to fetch APOD: %v", err)
	}
	log.Println("APOD:", apod.Date, apod.Title, apod.MediaType)
	log.Println("URL:", apod.URL)

	caption := apod.Title + "\n" + apod.Date

	failed := 0
	for _, chatID := range cfg.ChatIDs {
		if apod.MediaType == "video" {
			err = telegram.SendMessage(cfg.BotToken, chatID, caption+"\n🎬 Сегодня видео: "+apod.URL)
		} else {
			err = telegram.SendPhoto(cfg.BotToken, chatID, apod.URL, caption)
		}
		if err != nil {
			log.Printf("Failed to send to %s: %v", chatID, err)
			failed++
			continue
		}
		log.Println("sent to", chatID)
	}

	if failed > 0 {
		log.Fatalf("%d of %d sends failed", failed, len(cfg.ChatIDs))
	}
}
