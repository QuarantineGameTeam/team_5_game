package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"team_5_game/config"
	"team_5_game/model/telegram"
)

func ProcessWebhookMessage(update *telegram.Update) {
	log.Println("Processing webhook message:", convertToString(update))
	message := update.Message
	callbackQuery := update.CallbackQuery

	if message != nil {
		if message.Text == "/start" && isCommand(message) {
			RegisterUser(message)
			return
		}
	}

	if callbackQuery != nil {
		switch {
		case strings.HasPrefix(callbackQuery.Data, "CLAN_SELECT"):
			ProcessClanSelection(callbackQuery)
		case strings.HasPrefix(callbackQuery.Data, "START_BATTLE"):
			ProcessBattleStarting(callbackQuery)
		case strings.HasPrefix(callbackQuery.Data, "PRESS"):
			re := regexp.MustCompile("[0-9]+")
			position, err := strconv.Atoi(re.FindAllString(callbackQuery.Data, -1)[0]) // getting number from callbackQuery.Data
			if err != nil {
				fmt.Println("Could not convert Data to int:", err)
			}
			emoji, clanEmoji, _ := ClanParameters(callbackQuery)
			
			sendHintIfUnavailable(callbackQuery, emoji)

			AppendUserTrack(callbackQuery, position)

			replyMarkup := SendBattlefield(position, emoji, clanEmoji, callbackQuery)                            // getting field markup
			EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, &replyMarkup) // editing previous markup
			IsFull(callbackQuery)
		}
	}
}

func sendHintIfUnavailable(callbackQuery *telegram.CallbackQuery, emoji string) {
	if strings.HasPrefix(callbackQuery.Data, "PRESS_UNAVAILABLE") {
		SendAnswerCallbackQuery(callbackQuery.ID, "☹️You can capture neighboring cells only:\n" + "↖️🔼↗️\n◀️" + emoji + "▶️\n↙️🔽↘️", true)
	}
}

func SendMessage(chatID int64, message string, replyMarkup *telegram.InlineKeyboardMarkup) {
	err := sendMessage(chatID, message, replyMarkup)
	if err != nil {
		log.Println("Error in sending message:", err)
	}
}

func EditMessageReplyMarkup(chatID int64, messageID int64, replyMarkup *telegram.InlineKeyboardMarkup) {
	err := editMessageReplyMarkup(chatID, messageID, replyMarkup)
	if err != nil {
		log.Println("Error in editing message reply markup:", err)
	}
}

func SendAnswerCallbackQuery(callbackQueryID string, text string, showAlert bool) {
	err := sendAnswerCallbackQuery(callbackQueryID, text, showAlert)
	if err != nil {
		log.Println("Error in sending alert:", err)
	}
}

func convertToString(update *telegram.Update) string {
	out, err := json.Marshal(update)
	if err != nil {
		log.Println("Could not marshal update message", err)
		return "[Unable to convert to string]"
	}

	return string(out)
}

func sendMessage(chatID int64, message string, replyMarkup *telegram.InlineKeyboardMarkup) error {
	log.Println("Sending message to the chat:", chatID, " message: ", message)
	reqBody := &telegram.NewMessage{
		ChatID:      chatID,
		Text:        message,
		ReplyMarkup: replyMarkup,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken()+"/sendMessage",
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

func editMessageReplyMarkup(chatID int64, messageID int64, replyMarkup *telegram.InlineKeyboardMarkup) error {
	log.Println("Editing message reply markup chat:", chatID, " message: ", messageID)
	reqBody := &telegram.EditMessageReplyMarkup{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: replyMarkup,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken()+"/editMessageReplyMarkup",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	log.Println("Message changed successfully")
	return nil
}

func sendAnswerCallbackQuery(callbackQueryID string, text string, showAlert bool) error {
	log.Println("Sending alert message: ", text)
	reqBody := &telegram.AnswerCallbackQuery{
		CallbackQueryID: callbackQueryID,
		Text:            text,
		ShowAlert:       showAlert,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken()+"/answerCallbackQuery",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	log.Println("Alert sent successfully")
	return nil
}

func isCommand(message *telegram.Message) bool {
	if message.Entities == nil || len(*message.Entities) == 0 {
		return false
	}

	entity := (*message.Entities)[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}
