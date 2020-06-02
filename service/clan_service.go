package service

import (
	"log"
	"team_5_game/model/telegram"
)

func StartClanSelection(message *telegram.Message) {
	log.Println("Start clan selection for user ID", message.From.ID)

	replyMarkup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Clan 1", "CLAN_SELECT_1"),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Clan 2", "CLAN_SELECT_2"),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Clan 3", "CLAN_SELECT_3"),
		),
	)

	SendMessage(message.Chat.ID, "Please select a clan", &replyMarkup)
}

func ProcessClanSelection(callbackQuery *telegram.CallbackQuery) {
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)

	SaveUserClan(callbackQuery)

	user, err := GetUserFromDB(callbackQuery.From.ID) 
	if err != nil {
		log.Println("Could not get user", err)
	}

	SendMessage(callbackQuery.Message.Chat.ID, "Welcome to " + user.Clan + " clan)", nil)
}
