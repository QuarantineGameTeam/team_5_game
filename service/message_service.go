package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"team_5_game/config"
	"team_5_game/model/telegram"
)

func ProcessWebhookMessage(update *telegram.Update) {
	log.Println("Processing webhook message:", convertToString(update))
	message := update.Message
	callbackQuery := update.CallbackQuery

	if message != nil && isCommand(message) {
		log.Println("Received command:", message.Text)
		switch message.Text {
		case "/start":
			err := RegisterUser(message)
			if err == nil {
				log.Println("Start clan selection for user ID", message.From.ID)
				SendClanSelectionMenu(message)
			}
		case "/help":
			sendHelp(message)
		case "/restart":
			RestartGame(message)
		case "/stat":
			sendStatistic(message)
		default:
			log.Println("Received undefined command:", message.Text)
		}
		return
	}

	if callbackQuery != nil {
		switch {
		case strings.HasPrefix(callbackQuery.Data, "CLAN_SELECT"):
			ProcessClanSelection(callbackQuery)
			SendStartBattleMessage(callbackQuery)
		case strings.HasPrefix(callbackQuery.Data, "START_BATTLE"):
			ProcessBattleStarting(callbackQuery)
		case strings.HasPrefix(callbackQuery.Data, "PRESS"):
			ProcessNextMove(callbackQuery)
		}
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

func sendHintIfUnavailable(callbackQuery *telegram.CallbackQuery, emoji string) {
	SendAnswerCallbackQuery(
		callbackQuery.ID,
		"☹️You can capture neighboring cells only:\n"+"↖️🔼↗️\n◀️"+emoji+"▶️\n↙️🔽↘️",
		true,
	)
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

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println("Couldn't read response body:", err)
	}
	log.Printf("%s\n", resBody)

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

func sendHelp(message *telegram.Message) {
	log.Println("Send response message to `/help`")
	user, err := GetUserFromDB(message.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}
	SendMessage(
		message.Chat.ID,
		"...\nЗначення символів:\n▪️ - незахоплена територія\n"+Clans[user.ClanID].PlayerSign+" - ваша позиція (клан "+Clans[user.ClanID].Name+")\n🔹 - територія, захоплена кланом Blue Jays\n🔻 - територія, захоплена кланом Cardinals\n🔸 - територія, захоплена кланом Golden Orioles\nДоступний список команд:\n/restart - розпочати новий бій\n/start - розпочати гру\n/stat - статистика гри",
		nil)
	log.Println("Response message to `/help` has been successfully sent")
}

func sendStatistic(message *telegram.Message) {
	SendMessage(message.Chat.ID, "Here will be game statistic", nil)
	log.Println("Response message to `/stat` has been successfully sent")
}
