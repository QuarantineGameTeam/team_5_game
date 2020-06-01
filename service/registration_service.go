package service

import (
	"log"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

func RegisterUser(message *telegram.Message) {
	log.Println("Start user registration")

	user, _ := GetUserFromDB(message.From.ID)
	if user != nil {
		SendMessage(
			message.Chat.ID,
			"Hello "+message.From.FirstName+" you already registered!!!",
			nil)
		log.Println("User already registered")
	} else {
		user := database.User{
			ID:            message.From.ID,
			FirstName:     message.From.FirstName,
			Clan:          &database.Clan{ID: 0, Name: "NO_CLAN"},
			BattleCounter: 0,
			WinCounter:    0,
		}
		err := SaveUserToDB(&user)
		if err == nil {
			SendMessage(
				message.Chat.ID,
				"Hello "+message.From.FirstName+" thank you for registration!!!",
				nil)
			log.Println("User successfully registered")

			StartClanSelection(message)
		} else {
			SendMessage(
				message.Chat.ID,
				"Something went wrong please contact the administrator!!!",
				nil)
			log.Println("User not registered")
		}
	}
}
