package chat

import (
	"encoding/json"
	"time"
)

var acceptCheatMessageTypes = map[string]string{
	"text": "text",
}

type chatMessagePayload struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func (c chatMessagePayload) Parse(data interface{}) (chatMessagePayload, error) {
	p := chatMessagePayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}
