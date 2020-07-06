package service

import (
	"log"
	"math/rand"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
	"time"
)

func CreateBot(callbackQuery *telegram.CallbackQuery) (*database.User, error) {
	log.Println("Creating the zombie player")
	var zombie database.User

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get User", err)
		return nil, err
	}

	//zombie makes clan choice
	zombie.ClanID = 1
	if zombie.ClanID == user.ClanID {
		zombie.ClanID = 2
	}

	rand.Seed(time.Now().UnixNano())
	zombie.ID = user.ID + 2
	zombie.CurrentPos = Clans[zombie.ClanID].StartPosition
	zombie.CurrentBattle = user.CurrentBattle

	err = SaveUserToDB(&zombie)
	if err != nil {
		log.Println("Could not create zombie player", err)
		return nil, err
	}
	log.Println("Zombie player succesfully created, ID", zombie.ID)
	return &zombie, nil
}

func AddBotToBattle(battle *database.Battle, zombie *database.User) {
	log.Println("Adding zombie to battle")
	ZombiePositionCapture(battle, zombie)
	log.Println("Zombie successfully added to battle")
}

func BotMakesNextStep(callbackQuery *telegram.CallbackQuery) {
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

	walkRadius := ZombieWalkRadius(battle, zombie.CurrentPos)

	for _, cell := range walkRadius {
		if cell > 0 && cell < 26 {
			xy := ToCoordinates(cell)
			if !battle.Sector[xy.Row][xy.Col].IsCaptured { // Free cell search
				cellChoice = cell
				break
			}
		}
	}

	// If there are no free cells (isCaptured == true)
	for isCaptured {
		cellChoice = walkRadius[rand.Intn(len(walkRadius)-1)] // Choise random cell
		if cellChoice > 0 && cellChoice < 26 {
			break
		}
	}

	zombie.CurrentPos = cellChoice
	log.Println("ZombiePlayer POSITION:", zombie.CurrentPos)
	SaveUserToDB(zombie)
	ZombiePositionCapture(battle, zombie)
}

func ZombiePositionCapture(battle *database.Battle, zombie *database.User) {
	log.Printf("Capturing position %d by zombie\n", zombie.CurrentPos)
	xy := ToCoordinates(zombie.CurrentPos)

	if !battle.Sector[xy.Row][xy.Col].IsCaptured {
		battle.Sector[xy.Row][xy.Col].OwnedBy[0] = zombie.ClanID
		battle.Sector[xy.Row][xy.Col].IsCaptured = true
		err := SaveBattleToDB(battle)
		if err != nil {
			log.Println("Could not capture position", err)
			return
		}
		log.Println("Position successfully captured")
		return
	}
	log.Println("Position is already captured by another clan")
}

func ZombieWalkRadius(battle *database.Battle, pos int) (availableSectors []int) {
	log.Println("Trying to determine available sectors for zombies")
	yx := ToCoordinates(pos)

	for row := yx.Row - 1; row <= yx.Row+1; row++ {
		if row >= FieldHeigh || row < 0 {
			continue
		}
		for col := yx.Col - 1; col <= yx.Col+1; col++ {
			if col >= FieldWidth || col < 0 {
				continue
			}
			availableSectors = append(availableSectors, battle.Sector[row][col].ID)
		}
	}
	log.Println("Available territory:", availableSectors)

	return availableSectors
}
