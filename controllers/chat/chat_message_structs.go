package chat

import (
	"encoding/json"
	"time"
)

var acceptCheatMessageTypes = map[string]string{
	"text": "text",
}

type chatMessageResponse struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        uint      `json:"to"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type sendChatMessageRequestPayload struct {
	From    string `json:"from"`
	To      uint   `json:"to"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (sendChatMessageRequestPayload) parse(data interface{}) (sendChatMessageRequestPayload, error) {
	p := sendChatMessageRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}

func (c sendChatMessageRequestPayload) check() (isNotValid bool) {
	if _, ok := acceptCheatMessageTypes[c.Type]; !ok {
		return true
	}

	return false
}

type getChatMessageRequestPayload struct {
	ChatRoomID uint `json:"chatRoomId"`
}

type chatMessagesResponse []chatMessageResponse

func (c getChatMessageRequestPayload) Parse(data interface{}) (getChatMessageRequestPayload, error) {
	p := getChatMessageRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}
