package service

import (
	"encoding/json"
	"log"
	"net/http"
	"team_5_game/model_telegram"
)

func updateHandler(_ http.ResponseWriter, req *http.Request) {
	log.Println("Received update message")

	body := &model_telegram.Update{}

	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Println("Could not decode request body", err)
		return
	}

	ProcessUpdateMessage(body)
}

func CreateHttpServer(port string) {
	log.Println("Starting HTTP server")

	http.HandleFunc("/update", updateHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln("Could not start server", err)
	}

	log.Println("HTTP server started successfully")
}
