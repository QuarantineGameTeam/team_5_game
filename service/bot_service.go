package service

import (
	"log"
)

func StartBot() {
	log.Println("Starting Telegram bot.")

	CreateHttpServer()
	RegisterWebhook()

	select {} // block forever
}
