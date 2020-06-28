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
	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	ZombieWalking(callbackQuery) //zombie takes a step

	// Getting number from callbackQuery.Data
	re := regexp.MustCompile("[0-9]+")
	position, err := strconv.Atoi(re.FindAllString(callbackQuery.Data, -1)[0])
	if err != nil {
		log.Println("Could not convert Data to int:", err)
	}

	sendHintIfUnavailable(callbackQuery, Clans[user.ClanID].PlayerSign)
	AppendUserTrack(callbackQuery, position)

	// Get the next field markup.
	replyMarkup := SendBattlefield(position, user.ClanID, callbackQuery)
	// Editing previous markup
	EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, &replyMarkup)
	IsFull(callbackQuery)
}

func ProcessBattleStarting(callbackQuery *telegram.CallbackQuery) {

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	battleID := user.ID
	if user.CurrentBattle == 0 {
		err = SetNextBattle(battleID)
		if err != nil {
			log.Println("Could not set next battlefield", err)
			return
		}
		err = SetUserCurrentBattle(user.ID, battleID)
		if err != nil {
			return
		}
		LaunchZombieUser(callbackQuery) //adding bot-player to game
	} else {
		log.Println("Could not start new battle: user is in battle now")
		return
	}

	SendMessage(callbackQuery.Message.Chat.ID, "Your emoji: "+Clans[user.ClanID].PlayerSign, nil)

	replyMarkup := SendBattlefield(Clans[user.ClanID].StartPosition, user.ClanID, callbackQuery)     // getting field markup
	SendMessage(callbackQuery.Message.Chat.ID, "Select the cell you want to capture:", &replyMarkup) // creating message with new markup

	AppendUserTrack(callbackQuery, Clans[user.ClanID].StartPosition)
}

func AppendUserTrack(callbackQuery *telegram.CallbackQuery, position int) {
	log.Println("Start track saving")

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		return
	}
	battle, err := GetBattleFromDB(callbackQuery.From.ID)
	if err != nil {
		return
	}

	battle.Sector[position-1].OwnedBy[0] = callbackQuery.From.ID
	battle.Sector[position-1].IsCaptured = true
	user.CurrentPos = position

	SaveUserToDB(user)
	SaveBattleToDB(battle)
}

func SendBattlefield(position int, clanID int, callbackQuery *telegram.CallbackQuery) telegram.InlineKeyboardMarkup {
	var unknownTerritory string
	unknownTerritory = "▪️"
	//user, _ := GetUserFromDB(callbackQuery.From.ID)
	battle, _ := GetBattleFromDB(callbackQuery.From.ID)

	replyMarkup := telegram.InlineKeyboardMarkup{}

	min := 1
	max := 5

	for i := 1; i <= 5; i++ {
		var row []telegram.InlineKeyboardButton

		for j := min; j <= max; j++ {
			var btn telegram.InlineKeyboardButton
			if j == position {
				btn = telegram.NewInlineKeyboardButtonData(Clans[clanID].PlayerSign, "PRESS_"+strconv.Itoa(j))
			} else if IsAvailable(j, position) {
				btn = telegram.NewInlineKeyboardButtonData(unknownTerritory, "PRESS_"+strconv.Itoa(j))
			} else {
				btn = telegram.NewInlineKeyboardButtonData(unknownTerritory, "PRESS_UNAVAILABLE_"+strconv.Itoa(position))
			}
			if battle.Sector[j-1].IsCaptured {
				if btn.Text == unknownTerritory && len(battle.Sector) > 1 {
					btn.Text = Clans[clanID].ClanSign
				} else {
					btn.Text += Clans[clanID].ClanSign
				}
			}
			row = append(row, btn)
		}

		min += 5
		max += 5

		replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, row)
	}

	return replyMarkup
}

func IsAvailable(j int, position int) bool {
	res := false
	for _, element := range AvailableTerritory(position) {
		if j == element {
			res = true
			break
		}
	}
	return res
}

func AvailableTerritory(position int) []int {
	var availableTerritory []int
	if !(position%5 == 1) {
		availableTerritory = append(availableTerritory, position-1)
		availableTerritory = append(availableTerritory, position+4)
		availableTerritory = append(availableTerritory, position-6)
	}
	if !(position%5 == 0) {
		availableTerritory = append(availableTerritory, position+1)
		availableTerritory = append(availableTerritory, position+6)
		availableTerritory = append(availableTerritory, position-4)
	}
	availableTerritory = append(availableTerritory, position+5)
	availableTerritory = append(availableTerritory, position-5)
	availableTerritory = append(availableTerritory, position)
	return availableTerritory
}

func IsFull(callbackQuery *telegram.CallbackQuery) {
	user, _ := GetUserFromDB(callbackQuery.From.ID)
	battle, _ := GetBattleFromDB(callbackQuery.From.ID)
	res := true
	for _, point := range battle.Sector {
		for _, capture := range point.OwnedBy {
			if capture != 0 {
				point.IsCaptured = true
			}
			if !point.IsCaptured {
				res = false
				break
			}
		}
	}

	if res {
		EditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, nil)
		SendMessage(callbackQuery.Message.Chat.ID, "Game over!", nil)
		SetUserCurrentBattle(user.ID, 0)
		SendStartBattleMessage(callbackQuery)
	}
}

func SetNextBattle(id int64) error {
	// TO DO: find out how to check existanse of the DB key without getting all record.
	battle, err := GetBattleFromDB(id)
	if battle == nil {
		RegisterBattle(id)
	} else if err == nil {
		resetExistingBattle(battle)
	} else {
		log.Println("Could not set next battle", err)
		return err
	}
	return nil
}

func SetUserCurrentBattle(userID int64, battleID int64) error {
	user, err := GetUserFromDB(userID)
	if err != nil {
		log.Println("Could not get user", err)
		return err
	}
	user.CurrentBattle = battleID

	err = SaveUserToDB(user)
	if err != nil {
		log.Println("Could not save user's current battle to DB", err)
		return err
	} else {
		log.Printf("User's current battle successfully saved to DB, ID: %d\n", user.CurrentBattle)
	}
	return nil
}

func resetExistingBattle(battle *database.Battle) {
	for i := range battle.Sector {
		// Reset all of player domain marks in i-th sector.
		for j := range battle.Sector[i].OwnedBy {
			battle.Sector[i].OwnedBy[j] = 0
		}
		battle.Sector[i].IsCaptured = false
	}
	err := SaveBattleToDB(battle)
	if err != nil {
		log.Println("Could not reset existing battle", err)
	}
}
