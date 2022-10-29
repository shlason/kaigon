package chat

import (
	"encoding/json"
	"time"
)

type chatMessagePayload struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Text      string    `json:"text"`
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
