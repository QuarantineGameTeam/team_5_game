package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"team_5_game/config"
	"team_5_game/model_telegram"
)

func ProcessUpdateMessage(update *model_telegram.Update) {
	log.Println("Processing update message:", updateMessageToString(update))

	switch update.Message.Text {
	case "/start":
		registerUser(update)
	default:
		log.Println("Not implemented")
	}

}

func updateMessageToString(update *model_telegram.Update) string {
	out, err := json.Marshal(update)
	if err != nil {
		log.Println("Could not marshal update message", err)
		return "[Unable to convert to string] "
	}

	return string(out)
}

func sendMessage(chatID int64, message string) error {
	reqBody := &model_telegram.NewMessage{
		ChatID: chatID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken+"/sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func registerUser(update *model_telegram.Update) {
	//TODO: register User
	if err := sendMessage(update.Message.Chat.ID, "Hello new user!!!"); err != nil {
		log.Println("Error in sending message:", err)
		return
	}
}
