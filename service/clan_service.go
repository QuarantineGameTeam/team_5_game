package service

import (
	"log"
	"strconv"
	"strings"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

// Clans map contains information about all clans (database model)
var Clans = map[int]database.Clan{
	1: database.Clan{"Blue Jays", "💙", "🔹", 20},
	2: database.Clan{"Golden Orioles", "🧡", "🔸", 3},
	3: database.Clan{"Cardinals", "❤️", "🔻", 21},
}

func StartClanSelection(message *telegram.Message) {
	log.Println("Start clan selection for user ID", message.From.ID)

	replyMarkup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(Clans[1].ClanSign+" "+Clans[1].Name, "CLAN_SELECT_1"),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(Clans[2].ClanSign+" "+Clans[2].Name, "CLAN_SELECT_2"),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(Clans[3].ClanSign+" "+Clans[3].Name, "CLAN_SELECT_3"),
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

	SendMessage(callbackQuery.Message.Chat.ID, "Welcome to "+Clans[user.ClanID].Name+" clan)", nil)

	SendStartBattleMessage(callbackQuery)
}

func SaveUserClan(callbackQuery *telegram.CallbackQuery) {
	log.Println("Start clan saving")

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	clanID := strings.Trim(string(callbackQuery.Data), "SELECT_CLAN_")
	user.ClanID, err = strconv.Atoi(clanID)
	if err != nil {
		log.Println("Could not get user's clan ID", err)
		return
	}

	err = SaveUserToDB(user)
	if err != nil {
		log.Println("Could not save clan to DB", err)
	} else {
		log.Printf("User's clan successfully saved to DB, ID: %d, Name: %s\n",
			user.ClanID, Clans[user.ClanID].Name)
	}
}
