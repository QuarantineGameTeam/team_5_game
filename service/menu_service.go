package service

import "team_5_game/model/telegram"

func СhooseАction(message *telegram.Message) {
	StartBattle(message)
	//TODO: implement another choice
}
