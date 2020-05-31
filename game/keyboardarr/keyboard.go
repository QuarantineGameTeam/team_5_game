package keyboardarr

import (
	"encoding/json"
	"fmt"
)

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url,omitempty"`                              // optional
	CallbackData                 string `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"` // optional
}

//Тело запроса с массивом кнопок
type SendMessageReqBody struct {
	ChatID      int64                `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

var MaxRows int = 5
var MaxColumns int = 5

//Нужно придумать как кнопки итерировать, чтобы задавать матрицы нужного размера
func NewNumericKeyboard() InlineKeyboardMarkup {
	var numericKeyboard InlineKeyboardMarkup

	numericKeyboard = NewInlineKeyboardMarkup(
		NewInlineKeyboardRow(
			NewInlineKeyboardButtonData(ButtonText1, ButtonData1),
			NewInlineKeyboardButtonData(ButtonText2, ButtonData2),
			NewInlineKeyboardButtonData(ButtonText3, ButtonData3),
			NewInlineKeyboardButtonData(ButtonText4, ButtonData4),
			NewInlineKeyboardButtonData(ButtonText5, ButtonData5),
		),

		NewInlineKeyboardRow(
			NewInlineKeyboardButtonData(ButtonText6, ButtonData6),
			NewInlineKeyboardButtonData(ButtonText7, ButtonData7),
			NewInlineKeyboardButtonData(ButtonText8, ButtonData8),
			NewInlineKeyboardButtonData(ButtonText9, ButtonData9),
			NewInlineKeyboardButtonData(ButtonText10, ButtonData10),
		),

		NewInlineKeyboardRow(
			NewInlineKeyboardButtonData(ButtonText11, ButtonData11),
			NewInlineKeyboardButtonData(ButtonText12, ButtonData12),
			NewInlineKeyboardButtonData(ButtonText13, ButtonData13),
			NewInlineKeyboardButtonData(ButtonText14, ButtonData14),
			NewInlineKeyboardButtonData(ButtonText15, ButtonData15),
		),

		NewInlineKeyboardRow(
			NewInlineKeyboardButtonData(ButtonText16, ButtonData16),
			NewInlineKeyboardButtonData(ButtonText17, ButtonData17),
			NewInlineKeyboardButtonData(ButtonText18, ButtonData18),
			NewInlineKeyboardButtonData(ButtonText19, ButtonData19),
			NewInlineKeyboardButtonData(ButtonText20, ButtonData20),
		),

		NewInlineKeyboardRow(
			NewInlineKeyboardButtonData(ButtonText21, ButtonData21),
			NewInlineKeyboardButtonData(ButtonText22, ButtonData22),
			NewInlineKeyboardButtonData(ButtonText23, ButtonData23),
			NewInlineKeyboardButtonData(ButtonText24, ButtonData24),
			NewInlineKeyboardButtonData(ButtonText25, ButtonData25),
		),
	)

	return numericKeyboard
}

func RequestBody(chatID int64, text string) []byte {
	// Create the request body struct
	reqBody := &SendMessageReqBody{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: NewNumericKeyboard(),
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Print(err)
	}

	return reqBytes
}

//Тут пока ничего
func ChangeButton(updateData string) {
	fmt.Println(updateData)
}

//Далее идут функции создающие ряды кнопок, взял тут: https://github.com/go-telegram-bot-api/telegram-bot-api
// NewInlineKeyboardMarkup creates a new inline keyboard.
func NewInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton

	keyboard = append(keyboard, rows...)

	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

// NewInlineKeyboardRow creates an inline keyboard row with buttons.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// NewInlineKeyboardButtonURL creates an inline keyboard button with text
// which goes to a URL.
func NewInlineKeyboardButtonURL(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		URL:  url,
	}
}

// NewInlineKeyboardButtonSwitch creates an inline keyboard button with
// text which allows the user to switch to a chat or return to a chat.
func NewInlineKeyboardButtonSwitch(text, sw string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:              text,
		SwitchInlineQuery: sw,
	}
}

// NewInlineKeyboardButtonData creates an inline keyboard button with text
// and data for a callback.
func NewInlineKeyboardButtonData(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}
}
