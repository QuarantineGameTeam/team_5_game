package service

import (
	"fmt"
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
	youAre := fmt.Sprintf("You are: %v", "🔵")
	SendMessage(callbackQuery.Message.Chat.ID, youAre, nil)
	SendBattlefield(callbackQuery)
}

func SendBattlefield(callbackQuery *telegram.CallbackQuery) {
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
	replyMarkup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("🔹", "PRESS_1"),
			telegram.NewInlineKeyboardButtonData("⚪️", "PRESS_2"),
			telegram.NewInlineKeyboardButtonData("🔹", "PRESS_3"),
			telegram.NewInlineKeyboardButtonData("⚪️", "PRESS_4"),
			telegram.NewInlineKeyboardButtonData("⚪️", "PRESS_5"),
		),

		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(" 🔹 ", "PRESS_6"),
			telegram.NewInlineKeyboardButtonData(" 🔹 ", "PRESS_7"),
			telegram.NewInlineKeyboardButtonData(" 🔹 ", "PRESS_8"),
			telegram.NewInlineKeyboardButtonData("🔺🔸", "PRESS_9"),
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_10"),
		),

		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_11"),
			telegram.NewInlineKeyboardButtonData(" 🔹 ", "PRESS_12"),
			telegram.NewInlineKeyboardButtonData(" 🔺 ", "PRESS_13"),
			telegram.NewInlineKeyboardButtonData(" 🔸 ", "PRESS_14"),
			telegram.NewInlineKeyboardButtonData(" 🔸 ", "PRESS_15"),
		),

		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_16"),
			telegram.NewInlineKeyboardButtonData("🔵🔺", "PRESS_17"),
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_18"),
			telegram.NewInlineKeyboardButtonData(" 🔺 ", "PRESS_19"),
			telegram.NewInlineKeyboardButtonData(" 🔸 ", "PRESS_20"),
		),

		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_21"),
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_22"),
			telegram.NewInlineKeyboardButtonData(" ⚪️ ", "PRESS_23"),
			telegram.NewInlineKeyboardButtonData(" 🔺 ", "PRESS_24"),
			telegram.NewInlineKeyboardButtonData(" 🔸 ", "PRESS_25"),
		),
	)

	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture:", &replyMarkup)
}
