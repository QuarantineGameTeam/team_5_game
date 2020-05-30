package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"team_5_game/model/telegram"
)

const updatePath = "/update"

func updateHandler(_ http.ResponseWriter, req *http.Request) {
	log.Println("Received update message")

	body := &telegram.Update{}

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		log.Println("Could not decode request body", err)
		return
	}

	ProcessUpdateMessage(body)
}

func CreateHttpServer() {
	log.Println("Starting HTTP server")

	http.HandleFunc(updatePath, updateHandler)

	go func() {
		err := http.ListenAndServe(os.Getenv("SERVER_PORT"), nil)
		if err != nil {
			log.Fatalln("Could not start server", err)
		}
	}()

	log.Println("HTTP server started successfully")
}

func RegisterWebhook() {
	registerWebhook, _ := strconv.ParseBool(os.Getenv("REGISTER_WEBHOOK"))
	if registerWebhook {
		log.Println("Registering Webhook")
		reqBytes := []byte(`{"url":"` + os.Getenv("SERVER_URL") + updatePath + `"}`)

		_, err := http.Post(
			"https://api.telegram.org/bot"+os.Getenv("BOT_TOKEN")+"/setWebhook",
			"application/json",
			bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Fatalln("Could not register Webhook", err)
		}

		log.Println("Webhook registered successfully")
	}
}
