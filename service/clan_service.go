package service

import (
	"log"
	"strconv"
	"strings"
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

	SendMessage(callbackQuery.Message.Chat.ID, "Welcome to "+user.Clan.Name+" clan)", nil)

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
	user.Clan.ID, err = strconv.Atoi(clanID)
	if err != nil {
		log.Println("Could not get user's clan ID", err)
		return
	}

	switch user.Clan.ID {
	case 1:
		user.Clan.Name = "Blue Jays"
	case 2:
		user.Clan.Name = "Golden Orioles"
	case 3:
		user.Clan.Name = "Cardinals"
	default:
		log.Println("Could not set the clan's name", err)
		return
	}

	err = SaveUserToDB(user)
	if err != nil {
		log.Println("Could not save clan to DB", err)
	} else {
		log.Printf("User's clan successfully saved to DB, ID: %d, Name: %s\n",
			user.Clan.ID, user.Clan.Name)
	}
}
