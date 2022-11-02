package chat

import (
	"fmt"
	"net/http"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

type getAllChatRoomResponse struct {
	Data []models.ChatRoom
}

func getAllChatRoomHandler(msg message) {
	var chatRoomMembers []models.ChatRoomMember

	result := models.ChatRoomMember{}.ReadAllByAccountUUID(msg.Self.AccountUUID, &chatRoomMembers)
	fmt.Println(result.Error, msg.Self.AccountUUID)
	// TODO: 沒有聊天室時
	if result.Error == gorm.ErrRecordNotFound {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    controllers.SuccessCode,
			StatusCode:    http.StatusOK,
			StatusMessage: controllers.SuccessMessage,
			Payload:       []models.ChatRoom{},
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
			Payload:       []models.ChatRoom{},
		}
		return
	}

	var chatRoomIds []interface{}
	var chatRooms []models.ChatRoom

	fmt.Println("chat room members: ", len(chatRoomMembers))

	for _, chatRoomMember := range chatRoomMembers {
		chatRoomIds = append(chatRoomIds, chatRoomMember.ChatRoomID)
	}

	result = models.ChatRoom{}.ReadByIDs(chatRoomIds, &chatRooms)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	fmt.Println(chatRooms, len(chatRooms))
}
