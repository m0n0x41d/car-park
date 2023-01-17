package main

import (
	"car-park/telegram-bot/clients/telegram_client"
	"car-park/telegram-bot/consumer/event_consumer"
	"car-park/telegram-bot/events/telegram"
	"car-park/telegram-bot/storage/files"
	"log"
	"os"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "filestorage"
	batchSize   = 100
)

func main() {
	tgClient := telegram_client.New(tgBotHost, mustToken())
	eventsProcessor := telegram.New(tgClient, files.New(storagePath))

	log.Printf("[INFO]: service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("[FATAL]service stopped because of: %s", err.Error())
	}

}

func mustToken() string {
	token := os.Getenv("TELEGRAM_APIKEY")
	if token == "" {
		panic("No token in ENV. Set it with TELEGRAM_APIKEY key")
	} else {
		return token
	}
}
