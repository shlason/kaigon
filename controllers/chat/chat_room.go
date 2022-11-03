package chat

import (
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
