package main

import (
	"flag"
	"log"

	tgclient "github.com/x-goto/golang-tg-bot-demo/clients/telegram"
	event_consumer "github.com/x-goto/golang-tg-bot-demo/consumer/event-consumer"
	"github.com/x-goto/golang-tg-bot-demo/events/telegram"
	"github.com/x-goto/golang-tg-bot-demo/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage/user_links"
	batchSize   = 100
)

func main() {
	eventsHandler := telegram.New(tgclient.New(tgBotHost, mustToken()), files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventsHandler, eventsHandler, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-api-key", "", "telegram bot api token")
	flag.Parse()

	if *token == "" {
		log.Fatal("empty telegram bot api token")
	}

	return *token
}
