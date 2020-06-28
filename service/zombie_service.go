package service

import (
	"log"
	"math/rand"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
	"time"
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
	zombie.CurrentPos = Clans[zombie.ClanID].StartPosition
	zombie.CurrentBattle = user.CurrentBattle

	SaveUserToDB(&zombie)
	ZombiePositionCapture(&zombie)
}

func ZombieWalking(callbackQuery *telegram.CallbackQuery) {
	var cellChoice int
	var isCaptured bool = true
	zombie, err := GetUserFromDB(callbackQuery.From.ID + 2)
	if err != nil {
		log.Println("Could not get ZombieUser", err)
		return
	}
	battle, err := GetBattleFromDB(zombie.CurrentBattle)
	if err != nil {
		log.Println("Could not get Battle for ZombiUser", err)
		return
	}

	walkRadius := ZombieWalkRadius(zombie.CurrentPos)

	for _, cell := range walkRadius {
		if cell > 0 && cell < 26 {
			isCaptured = battle.Sector[cell-1].IsCaptured
			if !isCaptured { // Free cell search
				cellChoice = cell
				break
			}
		}
	}

	// If there are no free cells (isCaptured == true)
	for isCaptured {
		rand.Seed(time.Now().UnixNano())
		cellChoice = walkRadius[rand.Intn(len(walkRadius)-1)] // Choise random cell
		if cellChoice > 0 && cellChoice < 26 {
			break
		}
	}

	zombie.CurrentPos = cellChoice
	log.Println("ZombiePlayer POSITION:", zombie.CurrentPos)
	SaveUserToDB(zombie)
	ZombiePositionCapture(zombie)
}

func ZombiePositionCapture(zombie *database.User) {
	battle, err := GetBattleFromDB(zombie.CurrentBattle)
	if err != nil {
		log.Println("Could not get Battle for ZombiUser", err)
		return
	}

	battle.Sector[zombie.CurrentPos-1].OwnedBy[1] = zombie.ID
	battle.Sector[zombie.CurrentPos-1].IsCaptured = true

	SaveBattleToDB(battle)
}

func ZombieWalkRadius(position int) []int {
	var availableTerritory []int
	switch {
	case position%5 == 0:
		availableTerritory = append(availableTerritory, position-1)
		availableTerritory = append(availableTerritory, position+4)
		availableTerritory = append(availableTerritory, position-6)
		break
	case position%5 == 1:
		availableTerritory = append(availableTerritory, position+1)
		availableTerritory = append(availableTerritory, position+6)
		availableTerritory = append(availableTerritory, position-4)
		break
	default:
		availableTerritory = append(availableTerritory, position-6)
		availableTerritory = append(availableTerritory, position-4)
		availableTerritory = append(availableTerritory, position+1)
		availableTerritory = append(availableTerritory, position-1)
		availableTerritory = append(availableTerritory, position+4)
		availableTerritory = append(availableTerritory, position+6)
	}

	availableTerritory = append(availableTerritory, position+5)
	availableTerritory = append(availableTerritory, position-5)

	return availableTerritory
}
