package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"team_5_game/config"
	"team_5_game/model/telegram"
)

const webhookPath = "/webhook"

func webhookHandler(_ http.ResponseWriter, req *http.Request) {
	log.Println("Received webhook message.")

	body := &telegram.Update{}

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		log.Println("Could not decode request body", err)
		return
	}

	ProcessWebhookMessage(body)
}

func CreateHttpServer() {
	log.Println("Starting HTTP server.")

	http.HandleFunc(webhookPath, webhookHandler)

	go func() {
		err := http.ListenAndServe(config.ServerPort(), nil)
		if err != nil {
			log.Fatalln("Could not start server", err)
		}
	}()

	log.Println("HTTP server started successfully.")
}

func RegisterWebhook() {
	if config.RegisterWebhook() {
		log.Println("Registering Webhook.")
		reqBytes := []byte(`{"url":"` + config.ServerURL() + webhookPath + `"}`)

		_, err := http.Post(
			"https://api.telegram.org/bot"+config.BotToken()+"/setWebhook",
			"application/json",
			bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Fatalln("Could not register Webhook", err)
		}

		log.Println("Webhook registered successfully.")
	}
}
