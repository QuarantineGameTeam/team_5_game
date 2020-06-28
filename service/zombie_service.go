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
		log.Println("Could not get User", err)
		return
	}

	//zombie makes clan choice
	zombie.ClanID = 1
	if zombie.ClanID == user.ClanID {
		zombie.ClanID = 2
	}

	zombie.ID = user.ID + 2
	zombie.CurrentPos = Clans[zombie.ClanID].StartPosition - 1
	zombie.CurrentBattle = user.CurrentBattle

	SaveUserToDB(&zombie)
	ZombiePositionCapture(&zombie)
}

func ZombieWalking(callbackQuery *telegram.CallbackQuery) {
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
	SaveUserToDB(zombie)
	ZombiePositionCapture(zombie)
}

func ZombiePositionCapture(zombie *database.User) {
	battle, err := GetBattleFromDB(zombie.CurrentBattle)
	if err != nil {
		log.Println("Could not get Battle for ZombiUser", err)
		return
	}

	battle.Sector[zombie.CurrentPos].OwnedBy[1] = zombie.ID
	battle.Sector[zombie.CurrentPos].IsCaptured = true

	SaveBattleToDB(battle)
}
