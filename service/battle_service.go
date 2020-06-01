package service

import (
	"log"
	"team_5_game/model/telegram"
)

func StartBattle(message *telegram.Message) {
	log.Println("Send battle request for user ID", message.From.ID)

	replyMarkup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Go!", "START_BATTLE"),
		),

		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("later", "SKIP"),
		),
	)

	SendMessage(message.Chat.ID, "Do you want to start a battle?", &replyMarkup)
}

func ProcessBattleStarting(callbackQuery *telegram.CallbackQuery) {
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
	//TODO: Implement saving choice for the user and start battle
}
