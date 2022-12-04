package chat

import (
	"encoding/json"
	"time"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/utils"
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
	EnabledNotification bool   `json:"enabledNotification"`
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

type createRoomRequestPayload struct {
	Type             string  `json:"type"`
	Name             *string `json:"name"`
	Avatar           *string `json:"avatar"`
	InvitedUserEmail *string `json:"invitedUserEmail"`
}

var chatRoomTypes = map[string]string{
	"personal": "personal",
	"group":    "group",
}

func (p createRoomRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
	if _, ok := chatRoomTypes[p.Type]; !ok {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	if p.Type == chatRoomTypes["personal"] {
		if p.InvitedUserEmail == nil || !utils.IsValidEmailAddress(*p.InvitedUserEmail) {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	if p.Type == chatRoomTypes["group"] {
		if p.Avatar == nil || *p.Avatar == "" || p.Name == nil || *p.Name == "" {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	return controllers.JSONResponse{}, false
}

type getRoomInviteCodeResponse struct {
	Code string `json:"code"`
}
