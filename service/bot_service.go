package service

import (
	"log"
	"team_5_game/config"
)

func StartBot() {
	log.Println("Starting Telegram bot")

	CreateHttpServer(config.ServerPort)
}
