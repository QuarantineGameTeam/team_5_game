package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"team_5_game/model/telegram"
)

func ProcessUpdateMessage(update *telegram.Update) {
	log.Println("Processing update message:", updateMessageToString(update))

	message := update.Message
	if message != nil {
		if message.Text == "/start" && containsMessageType(message.Entities, "bot_command") {
			registerUser(message)
			return
		}
	}
}

func updateMessageToString(update *telegram.Update) string {
	out, err := json.Marshal(update)
	if err != nil {
		log.Println("Could not marshal update message", err)
		return "[Unable to convert to string]"
	}

	return string(out)
}

func sendMessage(chatID int64, message string) error {
	log.Println("Sending message to the chat:", chatID, " message: ", message)
	reqBody := &telegram.NewMessage{
		ChatID: chatID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+os.Getenv("BOT_TOKEN")+"/sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	log.Println("Message sent successfully")
	return nil
}

func containsMessageType(messageEntities []*telegram.MessageEntity, messageType string) bool {
	for _, messageEntity := range messageEntities {
		if messageEntity.Type == messageType {
			return true
		}
	}
	return false
}
