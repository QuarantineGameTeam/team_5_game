package database

type User struct {
	ID            int64   `json:"id"`
	FirstName     string  `json:"first_name"`
	ClanID        int     `json:"clan"`
	BattleCounter int     `json:"battle_counter"`
	WinCounter    int     `json:"win_counter"`
	Track       [25]int   `json:"track"`
	//	Clan          *Clan  `json:"clan"`
}

type Clan struct {
	// ID            int    `json:"id"`
	Name          string `json:"name"`
	PlayerSign    string `json:"player_sign"`
	ClanSign      string `json:"clan_sign"`
	StartPosition int    `json:"resp"`
}
