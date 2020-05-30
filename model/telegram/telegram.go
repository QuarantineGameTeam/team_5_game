package telegram

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int64            `json:"message_id"`
	From      *User            `json:"from"`
	Date      int              `json:"date"`
	Chat      *Chat            `json:"chat"`
	Text      string           `json:"text"`
	Entities  []*MessageEntity `json:"entities"`
}

type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type MessageEntity struct {
	Type string `json:"type"`
}

type NewMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
