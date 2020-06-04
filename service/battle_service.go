package service

import (
	"log"
	"strconv"
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

	var clanSelected string
	var startPosition int

	clanSelected, startPosition = clanParameters(callbackQuery)

	SendMessage(callbackQuery.Message.Chat.ID, "Your emoji: "+clanSelected, nil)
	SendBattlefield(startPosition, clanSelected, callbackQuery)
}

func SendBattlefield(position int, clanEmoji string, callbackQuery *telegram.CallbackQuery) {
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)

	var unknownTerritory string
	unknownTerritory = "▪️"

	replyMarkup := telegram.InlineKeyboardMarkup{}

	min := 1
	max := 5

	for i := 1; i <= 5; i++ {
		var row []telegram.InlineKeyboardButton

		for j := min; j <= max; j++ {
			var btn telegram.InlineKeyboardButton
			if j == position {
				btn = telegram.NewInlineKeyboardButtonData(clanEmoji, "PRESS_"+strconv.Itoa(j))
			} else {
				btn = telegram.NewInlineKeyboardButtonData(unknownTerritory, "PRESS_"+strconv.Itoa(j))
			}
			row = append(row, btn)
		}

		min += 5
		max += 5

		replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, row)
	}

	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture:", &replyMarkup)
}

func clanParameters(callbackQuery *telegram.CallbackQuery) (string, int) {
	var emoji string
	var startPosition int

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
	}

	switch user.Clan {
	case "CLAN_SELECT_1":
		emoji = "💜"
		startPosition = 20
	case "CLAN_SELECT_2":
		emoji = "💚"
		startPosition = 3
	case "CLAN_SELECT_3":
		emoji = "💛"
		startPosition = 24
	}

	return emoji, startPosition
}
