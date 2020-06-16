package service

import (
	"log"
	"regexp"
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

// ProcessNextMove - player's intention to go to the next sector.
func ProcessNextMove(callbackQuery *telegram.CallbackQuery) {
	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	// Getting number from callbackQuery.Data
	re := regexp.MustCompile("[0-9]+")
	position, err := strconv.Atoi(re.FindAllString(callbackQuery.Data, -1)[0])
	if err != nil {
		log.Println("Could not convert Data to int:", err)
	}

	sendHintIfUnavailable(callbackQuery, Clans[user.ClanID].PlayerSign)
	AppendUserTrack(callbackQuery, position)

	// Get the next field markup.
	replyMarkup := SendBattlefield(position, Clans[user.ClanID].PlayerSign, Clans[user.ClanID].ClanSign, callbackQuery)
	// Editing previous markup
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, &replyMarkup)
	IsFull(callbackQuery)
}

func ProcessBattleStarting(callbackQuery *telegram.CallbackQuery) {

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	SendMessage(callbackQuery.Message.Chat.ID, "Your emoji: "+Clans[user.ClanID].PlayerSign, nil)

	replyMarkup := SendBattlefield(Clans[user.ClanID].StartPosition, Clans[user.ClanID].PlayerSign, Clans[user.ClanID].ClanSign, callbackQuery) // getting field markup
	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture:", &replyMarkup)                                            // creating message with new markup

	AppendUserTrack(callbackQuery, Clans[user.ClanID].StartPosition)
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

func IsFull(callbackQuery *telegram.CallbackQuery) {
	user, _ := GetUserFromDB(callbackQuery.From.ID)
	res := true
	for _, point := range user.Track {
		if point == 0 {
			res = false
			break
		}
	}

	if res {
		EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
		SendMessage(callbackQuery.Message.Chat.ID, "Game over!", nil)
		ClearUserTrack(callbackQuery)
		SendStartBattleMessage(callbackQuery)
	}
}
