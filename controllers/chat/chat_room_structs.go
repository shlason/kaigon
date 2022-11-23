package chat

import (
	"encoding/json"
	"time"
)

type chatRoomMemberResponse struct {
	AccountUUID         string
	Name                string
	Avatar              string
	Theme               string
	EnabledNotification bool
	LastSeenAt          time.Time
}

type chatRoomInfoResponse struct {
	ID               uint                     `json:"id"`
	Type             string                   `json:"type"`
	MaximumMemberNum int                      `json:"maximumMemberNum"`
	Emoji            string                   `json:"emoji"`
	Name             string                   `json:"name"`
	Avatar           string                   `json:"avatar"`
	Members          []chatRoomMemberResponse `json:"members"`
}

type getAllChatRoomResponse []chatRoomInfoResponse

type updateChatRoomSettingRequestPayload struct {
	ChatRoomID uint   `json:"chatRoomId"`
	Emoji      string `json:"emoji"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

func (updateChatRoomSettingRequestPayload) parse(data interface{}) (updateChatRoomSettingRequestPayload, error) {
	p := updateChatRoomSettingRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}

type updateChatRoomSettingResponse struct {
	ChatRoomID uint   `json:"chatRoomId"`
	Emoji      string `json:"emoji"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

type updateChatRoomCustomSettingRequestPayload struct {
	ChatRoomID          uint   `json:"chatRoomId"`
	Theme               string `json:"theme"`
	EnabledNotification string `json:"enabledNotification"`
}

func (updateChatRoomCustomSettingRequestPayload) parse(data interface{}) (updateChatRoomCustomSettingRequestPayload, error) {
	p := updateChatRoomCustomSettingRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}

type updateChatRoomLastSeenRequestPayload struct {
	ChatRoomID uint `json:"chatRoomId"`
}

func (updateChatRoomLastSeenRequestPayload) parse(data interface{}) (updateChatRoomLastSeenRequestPayload, error) {
	p := updateChatRoomLastSeenRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}

type updateChatRoomCustomSettingResponse struct {
	LastSeenAt time.Time `json:"lastSeenAt"`
}
