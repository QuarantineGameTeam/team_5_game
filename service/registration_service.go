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
			"Hello "+message.From.FirstName+" you're already registered!!!",
			nil)
		log.Println("User is already registered")
	} else {
		user := database.User{
			ID:            message.From.ID,
			FirstName:     message.From.FirstName,
			ClanID:        0,
			BattleCounter: 0,
			WinCounter:    0,
			// Track:         nil,
			// Clan:          &database.Clan{ID: 0, Name: "NO_CLAN"},
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

func RegisterBattle(id int64) {
	log.Println("Start battle registration")

	battle, _ := GetBattleFromDB(id)
	if battle != nil {
		log.Println("Battle is already registered")
	} else {
		battle := database.Battle{
			ID: id,
		}
		for i := range battle.Sector {
			battle.Sector[i].ID = i + 1
		}

		err := SaveBattleToDB(&battle)
		if err != nil {
			log.Println("Battle not registered")
		} else {
			log.Println("Battle successfully registered")
		}
	}
}
