package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

type chatRoomMemberResponse struct {
	AccountUUID         string
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

func getAllChatRoomHandler(msg message) {
	// 使用者所擁有的聊天室列表
	var availableChatRooms []models.ChatRoomMember

	result := models.ChatRoomMember{}.ReadAllByAccountUUID(msg.Self.AccountUUID, &availableChatRooms)
	fmt.Println(result.Error, msg.Self.AccountUUID)
	// TODO: 沒有聊天室時
	if result.Error == gorm.ErrRecordNotFound {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.SuccessCode,
			StatusCode:    http.StatusOK,
			StatusMessage: controllers.SuccessMessage,
			Payload:       getAllChatRoomResponse{},
		}
		return
	}
	// TODO: 發生意外錯誤時
	if result.Error != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeServerDatabaseQueryGotError,
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: result.Error,
			Payload:       getAllChatRoomResponse{},
		}
		return
	}

	var availableChatRoomIds []interface{}
	var chatRoomInfoList []models.ChatRoom

	for _, chatRoomMember := range availableChatRooms {
		availableChatRoomIds = append(availableChatRoomIds, chatRoomMember.ChatRoomID)
	}

	result = models.ChatRoom{}.ReadAllByIDs(availableChatRoomIds, &chatRoomInfoList)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	var chatRoomMembers []models.ChatRoomMember

	result = models.ChatRoomMember{}.ReadAllByChatRoomIDs(availableChatRoomIds, &chatRoomMembers)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	var chatRoomMembersMapping = make(map[uint][]chatRoomMemberResponse)

	for _, chatRoomMember := range chatRoomMembers {
		chatRoomMembersMapping[chatRoomMember.ChatRoomID] = append(
			chatRoomMembersMapping[chatRoomMember.ChatRoomID],
			chatRoomMemberResponse{
				AccountUUID:         chatRoomMember.AccountUUID,
				Theme:               chatRoomMember.Theme,
				EnabledNotification: chatRoomMember.EnabledNotification,
				LastSeenAt:          chatRoomMember.LastSeenAt,
			},
		)
	}

	var response getAllChatRoomResponse

	for _, chatRoomInfo := range chatRoomInfoList {
		response = append(response, chatRoomInfoResponse{
			ID:               chatRoomInfo.ID,
			Type:             chatRoomInfo.Type,
			MaximumMemberNum: chatRoomInfo.MaximumMemberNum,
			Emoji:            chatRoomInfo.Emoji,
			Name:             chatRoomInfo.Name,
			Avatar:           chatRoomInfo.Avatar,
			Members:          chatRoomMembersMapping[chatRoomInfo.ID],
		})
	}
	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           msg.Cmd,
		CustomCode:    controllers.SuccessCode,
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       response,
	}
}

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

func updateChatRoomSettingHandler(clients map[string]client, msg message) {
	requestPayload, err := updateChatRoomSettingRequestPayload{}.parse(msg.Payload)

	if err != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeRequestPayloadFieldNotValid,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: controllers.ErrMessageRequestPayloadFieldNotValid,
			Payload:       nil,
		}
		return
	}

	m := controllers.GetFilteredNilRequestPayloadMap(&requestPayload)

	chatRoomModel := models.ChatRoom{
		Model: gorm.Model{
			ID: requestPayload.ChatRoomID,
		},
	}

	result := chatRoomModel.UpdateByID(m)

	if result.Error != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeServerDatabaseUpdateGotError,
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: result.Error,
			Payload:       nil,
		}
		return
	}

	var chatRoomMembers []models.ChatRoomMember

	result = models.ChatRoomMember{}.ReadAllByChatRoomID(requestPayload.ChatRoomID, &chatRoomMembers)

	if result.Error != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeServerDatabaseQueryGotError,
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: result.Error,
			Payload:       nil,
		}
		return
	}

	for _, chatRoomMember := range chatRoomMembers {
		toCli, ok := clients[chatRoomMember.AccountUUID]
		// TODO: 接收方不在線上時的處理
		if !ok {
			fmt.Printf("Friend: %s offline\n", chatRoomMember.AccountUUID)
			continue
		}

		toCli <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			StatusCode:    http.StatusOK,
			StatusMessage: controllers.SuccessMessage,
			Payload: updateChatRoomSettingResponse{
				ChatRoomID: requestPayload.ChatRoomID,
				Emoji:      chatRoomModel.Emoji,
				Name:       chatRoomModel.Name,
				Avatar:     chatRoomModel.Avatar,
			},
		}
	}
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

func updateChatRoomCustomSettingHandler(msg message) {
	requestPayload, err := updateChatRoomCustomSettingRequestPayload{}.parse(msg.Payload)

	if err != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeRequestPayloadFieldNotValid,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: controllers.ErrMessageRequestPayloadFieldNotValid,
			Payload:       nil,
		}
		return
	}

	m := controllers.GetFilteredNilRequestPayloadMap(&requestPayload)

	chatRoomMemberModel := models.ChatRoomMember{
		ChatRoomID:  requestPayload.ChatRoomID,
		AccountUUID: msg.Self.AccountUUID,
	}

	result := chatRoomMemberModel.UpdateByChatRoomIDAndAccountUUID(m)

	if result.Error != nil {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.ErrCodeServerDatabaseUpdateGotError,
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: result.Error,
			Payload:       nil,
		}
		return
	}

	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           msg.Cmd,
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}
}
