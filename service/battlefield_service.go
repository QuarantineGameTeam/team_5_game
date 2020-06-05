package service

import (
	"team_5_game/model/telegram"
)

type Battlefield struct {
	Points  *[25]Point `json:"points"`
	Players *[]Player  `json:"players"`
}

type Point struct {
	Text        string `json:"text"`
	TextForUser string `json:"text_for_user"`
	Number      int    `json:"number"`
}

type Player struct {
	User  *telegram.User `json:"user"`
	Point *Point         `json:"point"`
	Clan  *Clan          `json:"clan"`
}

type Clan struct {
	Sign       string `json:"sign"`
	PlayerSign string `json:"player_sign"`
}

func IsFull(battleField Battlefield) bool {
	res := true
	for _, point := range battleField.Points {
		if len(point.Text) == 0 {
			res = false
			break
		}
	}
	return res
}
func CapturePoint(point *Point, player *Player) {
	point.Text += player.Clan.Sign
	player.Point = point
}
