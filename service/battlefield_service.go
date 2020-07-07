package service

import (
	"log"
	"strconv"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

const (
	// Buttlefield dimensions
	FieldHeigh = 5
	FieldWidth = 5
)

type Coordinates struct {
	Row int
	Col int
}

func SaveUserPosition(user *database.User, battle *database.Battle, position int) {
	log.Println("Trying update the user's current position to", position)
	xy := ToCoordinates(position)

	if !battle.Sector[xy.Row][xy.Col].IsCaptured {
		battle.Sector[xy.Row][xy.Col].OwnedBy[0] = user.ClanID
		battle.Sector[xy.Row][xy.Col].IsCaptured = true
	}
	user.CurrentPos = position

	SaveUserToDB(user)
	SaveBattleToDB(battle)
	log.Println("User's current position successfully updated to", user.CurrentPos)
}

func SendBattlefield(user *database.User, battle *database.Battle, callbackQuery *telegram.CallbackQuery) telegram.InlineKeyboardMarkup {
	var unknownTerritory string = "▪️"
	replyMarkup := telegram.InlineKeyboardMarkup{}

	for i := range battle.Sector {
		var row []telegram.InlineKeyboardButton
		for j := range battle.Sector[i] {
			var btn telegram.InlineKeyboardButton
			if battle.Sector[i][j].ID == user.CurrentPos {
				// Player in this sector; print player's emoji.
				btn = telegram.NewInlineKeyboardButtonData(
					Clans[user.ClanID].PlayerSign,
					"PRESS_"+strconv.Itoa(battle.Sector[i][j].ID),
				)
			} else if battle.Sector[i][j].IsCaptured {
				// Print emoji of the sector owner clan
				btn = telegram.NewInlineKeyboardButtonData(
					generateOwnersSign(battle, i, j),
					//generateSectorOwnersEmoji(&battle.Sector[i]),
					"PRESS_"+strconv.Itoa(battle.Sector[i][j].ID),
				)
			} else {
				// Print unknownTerritory emoji
				btn = telegram.NewInlineKeyboardButtonData(
					unknownTerritory,
					"PRESS_"+strconv.Itoa(battle.Sector[i][j].ID),
				)
			}
			row = append(row, btn)
		}
		replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, row)
	}
	return replyMarkup
}

func IsUnavailableTerritory(user *database.User, position int) bool {
	// Determine coordinates of the user current position
	currentXY := ToCoordinates(user.CurrentPos)
	// Determine coordinates of the next step position
	newXY := ToCoordinates(position)
	// Check the difference between current and next position
	diff := currentXY.Col - newXY.Col
	if diff > 1 || diff < -1 {
		return true
	}
	diff = currentXY.Row - newXY.Row
	if diff > 1 || diff < -1 {
		return true
	}
	return false
}

// ToCoords converts map sector number into map coordinates.
func ToCoordinates(n int) (c Coordinates) {
	c.Row = (n - 1) / FieldWidth
	c.Col = (n - 1) % FieldWidth
	return c
}

func IsFull(callbackQuery *telegram.CallbackQuery) bool {
	battle, _ := GetBattleFromDB(callbackQuery.From.ID)
	for row := range battle.Sector {
		for col := range battle.Sector[row] {
			if battle.Sector[row][col].IsCaptured == false {
				return false
			}
		}
	}
	return true
}

func generateOwnersSign(battle *database.Battle, i int, j int) (pic string) {
	for _, v := range battle.Sector[i][j].OwnedBy {
		if v != 0 {
			pic += Clans[v].ClanSign
		}
	}
	return pic
}
