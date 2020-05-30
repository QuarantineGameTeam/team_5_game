package service

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	// loads values from config.env into the system
	err := godotenv.Load("config.env")
	if err != nil {
		log.Print("No config.env file found")
	}
}

func StartBot() {
	log.Println("Starting Telegram bot")

	CreateHttpServer()
	RegisterWebhook()

	select {} // block forever
}
