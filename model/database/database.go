package database

type User struct {
	ID                          int64  `json:"id"`
	FirstName                   string `json:"first_name"`
	ClanID                      int    `json:"clan"`
	BattleCounter               int    `json:"battle_counter"`
	WinCounter                  int    `json:"win_counter"`
	CurrentBattle               int64  `json:"current_battle"`
	CurrentBattlefieldMessageID int64  `json:"current_battlefield_message_ID"`
	CurrentPos                  int    `json:"current_battle_position"`
	//	Clan          *Clan  `json:"clan"`
}

type Clan struct {
	// ID            int    `json:"id"`
	Name          string `json:"name"`
	PlayerSign    string `json:"player_sign"`
	ClanSign      string `json:"clan_sign"`
	StartPosition int    `json:"resp"`
}

type Battle struct {
	ID     int64 `json:"id"`
	Sector [5][5]struct {
		ID         int    `json:"id"`
		OwnedBy    [2]int `json:"owners"`
		IsCaptured bool   `json:"captured"`
	}
}
