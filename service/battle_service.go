package service

import (
	"log"
	"regexp"
	"strconv"
	"team_5_game/model/database"
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
	log.Println("Start processing next move for", callbackQuery.Data)

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}
	battle, err := GetBattleFromDB(user.CurrentBattle)
	if err != nil {
		log.Println("Could not get battle", err)
		return
	}

	// Get number of next sector from callbackQuery.Data
	re := regexp.MustCompile("[0-9]+")
	nextStep, err := strconv.Atoi(re.FindAllString(callbackQuery.Data, -1)[0])
	if err != nil {
		log.Println("Could not convert Data to int:", err)
		return
	}
	// Check whether a next sector is unavailable for the next step
	if IsUnavailableTerritory(user, nextStep) {
		sendHintIfUnavailable(callbackQuery, Clans[user.ClanID].PlayerSign)
		return
	}

	SaveUserPosition(user, battle, nextStep)
	replyMarkup := SendBattlefield(user, battle, callbackQuery)
	// Editing previous markup
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, &replyMarkup)
	if IsFull(callbackQuery) {
		log.Println("Battlefield is full")
		EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
		SendMessage(callbackQuery.Message.Chat.ID, "Game over!", nil)
		SetUserCurrentBattle(user, 0)
		SendStartBattleMessage(callbackQuery)
	}
}

func ProcessBattleStarting(callbackQuery *telegram.CallbackQuery) {
	log.Println("Initializing new battle")

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	// Generate battle ID, create new battle and assign the user to it.
	battleID := user.ID
	if user.CurrentBattle == 0 {
		err = SetNextBattle(battleID)
		if err != nil {
			log.Println("Could not set next battlefield", err)
			return
		}
		err = SetUserCurrentBattle(user, battleID)
		if err != nil {
			return
		}
	} else {
		log.Println("Could not start new battle: user is in battle now")
		return
	}
	battle, err := GetBattleFromDB(battleID)
	if err != nil {
		log.Println("Could not get battle", err)
		return
	}
	// Initialize start position of the user.
	SaveUserPosition(user, battle, Clans[user.ClanID].StartPosition)
	// Send user the new battlefield.
	SendMessage(callbackQuery.Message.Chat.ID, "Your emoji: "+Clans[user.ClanID].PlayerSign, nil)
	replyMarkup := SendBattlefield(user, battle, callbackQuery)
	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture: ", &replyMarkup)
	log.Println("New battle is successfully initialized")
}

func SetNextBattle(id int64) error {
	// TO DO: find out how to check existanse of the DB key without getting all record.
	battle, err := GetBattleFromDB(id)
	if battle == nil {
		RegisterBattle(id)
	} else if err == nil {
		ResetExistingBattle(battle)
	} else {
		log.Println("Could not set next battle", err)
		return err
	}
	return nil
}

func SetUserCurrentBattle(user *database.User, battleID int64) error {
	log.Println("Setting user's current battle to", battleID)
	user.CurrentBattle = battleID

	err := SaveUserToDB(user)
	if err != nil {
		log.Println("Could not save user's current battle to DB", err)
		return err
	}
	log.Printf("User's current battle successfully saved to DB, ID: %d\n", user.CurrentBattle)
	return nil
}

func ResetExistingBattle(battle *database.Battle) {
	log.Println("Resetting the battle ID:", battle.ID)
	for row := range battle.Sector {
		for col := range battle.Sector[row] {
			// Reset ownership of the sector for both possible owner clans
			for i := range battle.Sector[row][col].OwnedBy {
				battle.Sector[row][col].OwnedBy[i] = 0
			}
			battle.Sector[row][col].IsCaptured = false
		}
	}
	err := SaveBattleToDB(battle)
	if err != nil {
		log.Println("Could not reset existing battle", err)
	}
}
