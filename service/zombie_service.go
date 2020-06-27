package service

import (
	"log"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

func LaunchZombieUser(callbackQuery *telegram.CallbackQuery) {
	var zombie database.User

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get ZombieUser", err)
		return
	}

	zombie.ID = user.ID + 2
	zombie.ClanID = 1
	zombie.CurrentPos = Clans[zombie.ClanID].StartPosition - 1
	zombie.CurrentBattle = user.CurrentBattle

	battle, err := GetBattleFromDB(zombie.CurrentBattle)
	if err != nil {
		return
	}

	battle.Sector[zombie.CurrentPos].OwnedBy[1] = zombie.ID
	battle.Sector[zombie.CurrentPos].IsCaptured = true

	SaveUserToDB(&zombie)
	SaveBattleToDB(battle)
}

func ZomnbieWalking(callbackQuery *telegram.CallbackQuery) {
	zombie, err := GetUserFromDB(callbackQuery.From.ID + 2)
	if err != nil {
		log.Println("Could not get ZombieUser", err)
		return
	}

	if zombie.CurrentPos == 24 {
		zombie.CurrentPos = 0
	} else {
		zombie.CurrentPos++
	}

	battle, err := GetBattleFromDB(zombie.CurrentBattle)
	if err != nil {
		return
	}

	battle.Sector[zombie.CurrentPos].OwnedBy[1] = zombie.ID
	battle.Sector[zombie.CurrentPos].IsCaptured = true

	SaveUserToDB(zombie)
	SaveBattleToDB(battle)
}
