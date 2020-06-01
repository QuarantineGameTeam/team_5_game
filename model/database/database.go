package database

type User struct {
	ID            int64  `json:"id"`
	FirstName     string `json:"first_name"`
	Clan          *Clan  `json:"clan"`
	BattleCounter int    `json:"battle_counter"`
	WinCounter    int    `json:"win_counter"`
}

type Clan struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
