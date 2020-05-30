package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"team_5_game/config"
	"team_5_game/model_telegram"
)

const updatePath = "/update"

func updateHandler(_ http.ResponseWriter, req *http.Request) {
	log.Println("Received update message")

	body := &model_telegram.Update{}

	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Println("Could not decode request body", err)
		return
	}

	ProcessUpdateMessage(body)
}

func CreateHttpServer() {
	log.Println("Starting HTTP server")

	http.HandleFunc(updatePath, updateHandler)

	go func() {
		if err := http.ListenAndServe(config.ServerPort, nil); err != nil {
			log.Fatalln("Could not start server", err)
		}
	}()

	log.Println("HTTP server started successfully")
}

func RegisterWebhook() {
	if config.RegisterWebhook == true {
		log.Println("Registering Webhook")
		reqBytes := []byte(`{"url":"` + config.ServerUrl + updatePath + `"}`)

		_, err := http.Post(
			"https://api.telegram.org/bot"+config.BotToken+"/setWebhook",
			"application/json",
			bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Fatalln("Could not register Webhook", err)
		}

		log.Println("Webhook registered successfully")
	}
}
