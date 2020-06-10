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
	var clanSelected string
	var startPosition int
	var clanEmoji string

	clanSelected, clanEmoji, startPosition = ClanParameters(callbackQuery)

	SendMessage(callbackQuery.Message.Chat.ID, "Your emoji: "+clanSelected, nil)

	replyMarkup := SendBattlefield(startPosition, clanSelected, clanEmoji, callbackQuery)            // getting field markup
	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture:", &replyMarkup) // creating message with new markup

	AppendUserTrack(callbackQuery, startPosition)
}

func SendBattlefield(position int, emoji string, clanEmoji string, callbackQuery *telegram.CallbackQuery) telegram.InlineKeyboardMarkup {
	var unknownTerritory string
	unknownTerritory = "▪️"
	user, _ := GetUserFromDB(callbackQuery.From.ID)

	replyMarkup := telegram.InlineKeyboardMarkup{}

	min := 1
	max := 5

	for i := 1; i <= 5; i++ {
		var row []telegram.InlineKeyboardButton

		for j := min; j <= max; j++ {
			var btn telegram.InlineKeyboardButton
			if j == position {
				btn = telegram.NewInlineKeyboardButtonData(emoji, "PRESS_"+strconv.Itoa(j))
			} else if IsAvailable(j, position) {
				btn = telegram.NewInlineKeyboardButtonData(unknownTerritory, "PRESS_"+strconv.Itoa(j))
			} else {
				btn = telegram.NewInlineKeyboardButtonData(unknownTerritory, "PRESS_UNAVAILABLE_"+strconv.Itoa(position))
			}
			if IsThere(j, user.Track) {
				if btn.Text == unknownTerritory && len(user.Track) > 1 {
					btn.Text = clanEmoji
				} else {
					btn.Text += clanEmoji
				}
			}
			row = append(row, btn)
		}

		min += 5
		max += 5

		replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, row)
	}

	return replyMarkup
}

func ClanParameters(callbackQuery *telegram.CallbackQuery) (string, string, int) {
	var emoji string
	var clanEmoji string
	var startPosition int

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
	}

	switch user.Clan {
	case "CLAN_SELECT_1":
		emoji = "💙"
		clanEmoji = "🔹"
		startPosition = 20
	case "CLAN_SELECT_2":
		emoji = "🧡"
		clanEmoji = "🔸"
		startPosition = 3
	case "CLAN_SELECT_3":
		emoji = "❤️"
		clanEmoji = "🔻"
		startPosition = 24
	}

	return emoji, clanEmoji, startPosition
}

func IsAvailable(j int, position int) bool {
	res := false
	for _, element := range AvailableTerritory(position) {
		if j == element {
			res = true
			break
		}
	}
	return res
}

func AvailableTerritory(position int) []int {
	var availableTerritory []int
	if !(position%5 == 1) {
		availableTerritory = append(availableTerritory, position-1)
		availableTerritory = append(availableTerritory, position+4)
		availableTerritory = append(availableTerritory, position-6)
	}
	if !(position%5 == 0) {
		availableTerritory = append(availableTerritory, position+1)
		availableTerritory = append(availableTerritory, position+6)
		availableTerritory = append(availableTerritory, position-4)
	}
	availableTerritory = append(availableTerritory, position+5)
	availableTerritory = append(availableTerritory, position-5)
	availableTerritory = append(availableTerritory, position)
	return availableTerritory
}

func IsThere(element int, arr [25]int) bool {
	res := false
	for _, elem := range arr {
		if elem == element {
			res = true
			break
		}
	}
	return res
}
