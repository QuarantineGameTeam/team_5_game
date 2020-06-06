package service

import (
	"log"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

func RegisterUser(message *telegram.Message) {
	log.Println("Start user registration.")

	user, _ := GetUserFromDB(message.From.ID)
	if user != nil {
		SendMessage(
			message.Chat.ID,
			"Привіт, "+message.From.FirstName+", із поверненням! 😊",
			nil)
		log.Println("User is already registered.")
	} else {
		user := database.User{
			ID:        message.From.ID,
			FirstName: message.From.FirstName,
			Clan:      "NO_CLAN",
			// Clan:          &database.Clan{ID: 0, Name: "NO_CLAN"},
			BattleCounter: 0,
			WinCounter:    0,
		}
		err := SaveUserToDB(&user)
		if err == nil {
			SendMessage(
				message.Chat.ID,
				"Привіт, "+message.From.FirstName+", ми на тебе чекали) Що ж, розпочнемо?",
				nil)
			log.Println("User successfully registered.")

			StartClanSelection(message)
		} else {
			SendMessage(
				message.Chat.ID,
				"Упс! Щось пішло не так - звернись до адміністратора.",
				nil)
			log.Println("User was not registered.")
		}

	}
}
