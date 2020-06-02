package service

import (
	"log"
	"team_5_game/model/telegram"
)

func SendStartBattleMessage(callbackQuery *telegram.CallbackQuery) {
	log.Println("Send battle request for user ID", callbackQuery.From.ID)

	replyMarkup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Go!", "START_BATTLE"),
		),
	)

	SendMessage(callbackQuery.Message.Chat.ID, "Do you want to start a battle?", &replyMarkup)
}

func ProcessBattleStarting(callbackQuery *telegram.CallbackQuery) {
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
	//TODO:  Implement saving choice for the user and start battle
}
