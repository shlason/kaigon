package chat

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

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

	// TODO: 聊天室成員的相關資訊 query 方式太耗資源，要優化或改變 DB schema 或想辦法用 cache
	var chatRoomMembers []models.ChatRoomMember
	var chatRoomMemberAccountUUIDs []interface{}
	var chatRoomMemberAccountSettings []models.AccountSetting
	var chatRoomMemberAccountProfiles []models.AccountProfile

	result = models.ChatRoomMember{}.ReadAllByChatRoomIDs(availableChatRoomIds, &chatRoomMembers)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	for _, chatRoomMember := range chatRoomMembers {
		chatRoomMemberAccountUUIDs = append(chatRoomMemberAccountUUIDs, chatRoomMember.AccountUUID)
	}

	result = models.AccountSetting{}.ReadAllByAccountUUIDs(chatRoomMemberAccountUUIDs, &chatRoomMemberAccountSettings)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	result = models.AccountProfile{}.ReadAllByAccountUUIDs(chatRoomMemberAccountUUIDs, &chatRoomMemberAccountProfiles)

	// TODO: 發生意外錯誤時
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	// TODO: 寫法需要優化
	var chatRoomMemberAccountSettingsMapping = make(map[string]models.AccountSetting)
	var chatRoomMemberAccountProfilesMapping = make(map[string]models.AccountProfile)
	var chatRoomMembersMapping = make(map[uint][]chatRoomMemberResponse)

	for _, chatRoomMemberAccountSetting := range chatRoomMemberAccountSettings {
		chatRoomMemberAccountSettingsMapping[chatRoomMemberAccountSetting.AccountUUID] = models.AccountSetting{
			AccountID:   chatRoomMemberAccountSetting.AccountID,
			AccountUUID: chatRoomMemberAccountSetting.AccountUUID,
			Name:        chatRoomMemberAccountSetting.Name,
			Locale:      chatRoomMemberAccountSetting.Locale,
		}
	}

	for _, chatRoomMemberAccountProfile := range chatRoomMemberAccountProfiles {
		chatRoomMemberAccountProfilesMapping[chatRoomMemberAccountProfile.AccountUUID] = models.AccountProfile{
			AccountID:   chatRoomMemberAccountProfile.AccountID,
			AccountUUID: chatRoomMemberAccountProfile.AccountUUID,
			Avatar:      chatRoomMemberAccountProfile.Avatar,
			Banner:      chatRoomMemberAccountProfile.Banner,
			Signature:   chatRoomMemberAccountProfile.Signature,
		}
	}

	for _, chatRoomMember := range chatRoomMembers {
		chatRoomMembersMapping[chatRoomMember.ChatRoomID] = append(
			chatRoomMembersMapping[chatRoomMember.ChatRoomID],
			chatRoomMemberResponse{
				AccountUUID:         chatRoomMember.AccountUUID,
				Name:                chatRoomMemberAccountSettingsMapping[chatRoomMember.AccountUUID].Name,
				Avatar:              chatRoomMemberAccountProfilesMapping[chatRoomMember.AccountUUID].Avatar,
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

func updateChatRoomLastSeenHandler(clients map[string]client, msg message) {
	requestPayload, err := updateChatRoomLastSeenRequestPayload{}.parse(msg.Payload)

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

	chatRoomMemberModel := models.ChatRoomMember{
		ChatRoomID:  requestPayload.ChatRoomID,
		AccountUUID: msg.Self.AccountUUID,
	}

	result := chatRoomMemberModel.UpdateByChatRoomIDAndAccountUUID(
		map[string]interface{}{"LastSeenAt": time.Now()},
	)

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
			Payload: updateChatRoomCustomSettingResponse{
				LastSeenAt: chatRoomMemberModel.LastSeenAt,
			},
		}
	}
}

func CreateRoom(c *gin.Context) {
	var requestPayload createRoomRequestPayload

	errResp, err := controllers.BindJSON(c, &requestPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	errResp, isNotValid := requestPayload.check()

	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	chatRoomModel := models.ChatRoom{
		Type:   requestPayload.Type,
		Avatar: requestPayload.Avatar,
		Name:   requestPayload.Name,
	}

	result := chatRoomModel.Create()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

func GetRoomInviteCode(c *gin.Context) {
	chatRoomID, err := strconv.ParseUint(c.Param("chatRoomID"), 10, 22)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	chatRoomMemberModel := models.ChatRoomMember{
		ChatRoomID:  uint(chatRoomID),
		AccountUUID: authPayload.AccountUUID,
	}
	result := chatRoomMemberModel.ReadByChatRoomIDAndAccountUUID()

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPermissionForbidden,
			Message: controllers.ErrMessageRequestPermissionForbidden,
			Data:    nil,
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	ChatRoomInviteCodeModel := models.ChatRoomInviteCode{
		ChatRoomID: chatRoomMemberModel.ChatRoomID,
	}
	err = ChatRoomInviteCodeModel.Create()

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisSetNXKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getRoomInviteCodeResponse{
			Code: ChatRoomInviteCodeModel.Code,
		},
	})
}
