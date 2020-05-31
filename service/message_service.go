package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"team_5_game/config"
	"team_5_game/game/keyboardarr"
	"team_5_game/model_telegram"
)

func ProcessUpdateMessage(update *model_telegram.Update) {
	log.Println("Processing update message:", updateMessageToString(update))

	switch update.Message.Text {
	case "/start":
		registerUser(update)
	case "/battle": //команда для показа клавиатуры
		newKeyboard(update)
	default:
		log.Println("Not implemented")
	}

}

// KeyboardUpd вызывает функцию, которая меняет подпись кнопки, выводит сообщение с информацией о выбранном секторе (пока ни с чем не связана)
func KeyboardUpd(callbackQuery *keyboardarr.Update) {
	if callbackQuery.CallbackQuery.Data != "" {
		keyboardarr.ChangeButton(callbackQuery.CallbackQuery.Data)
		sendMessage(callbackQuery.CallbackQuery.Message.Chat.ID, "Выбран сектор:")
		sendMessage(callbackQuery.CallbackQuery.Message.Chat.ID, callbackQuery.CallbackQuery.Data)
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

// добавил чтобы отправлять сообщения, по другой структуре из keyboardarr SendMessageReqBody
func clearSend(requestBytes []byte) error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken+"/sendMessage",
		"application/json",
		bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

//отправляет клавиатуру
func newKeyboard(update *model_telegram.Update) {
	chatID := update.Message.Chat.ID
	if err := clearSend(keyboardarr.RequestBody(chatID, "Выбери сектор:")); err != nil {
		log.Println("Error in sending message:", err)
		return
	}
}

func registerUser(update *model_telegram.Update) {
	//TODO: register User
	if err := sendMessage(update.Message.Chat.ID, "Hello new user!!!"); err != nil {
		log.Println("Error in sending message:", err)
		return
	}
}
