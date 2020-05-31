package connection

import (
	"team_5_game/model_telegram"
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./bot_db.db")
	if err != nil {
		log.Println("Could not open database", err)
	}

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL, clan_type INT, user_id INTEGER);")
	statement.Exec()

	return db
}

func CreateUser(message *model_telegram.Message) {
	db := OpenDatabase()
	var username string

	username = CreateUserName(message.From)
	idUser := message.From.ID

	if checkUser(idUser, db) {
		log.Println("User alredy exists")
	} else {
		statement, _ := db.Prepare("INSERT INTO users (username, clan_type, user_id) VALUES (?, ?, ?)")
		if _, err := statement.Exec(username, 1, idUser); err != nil {
			log.Println("Could not insert user:", err)
		} else { log.Println("User " + username + " insert to database") }
	}	
}

func checkUser(idUser int, db *sql.DB) bool {
	rows, _ := db.Query("SELECT id FROM users WHERE user_id = ?", idUser)
	var id int

	for rows.Next() {
		rows.Scan(&id)
	}

	if id == 0 {
		return false
	} else { 
		return true 
	}
}

func CreateUserName(user *model_telegram.User) string {
	if user.UserName != "" {
		return user.UserName
	} else if user.FirstName != "" {
		return user.FirstName + user.LastName
	} else {
		return user.LastName
	}
}

// type RegisteredUsers struct {
// 	ID       int    `json:"message_id"`
// 	UserID   int    `json:"user_id"`
// 	UserName string `json:"username"`
// 	ClanType int    `json:"clan_type"`
// }

// func getUser(db *sql.DB) RegisteredUsers{ // For returning user information

// }
