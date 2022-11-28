package chat

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

// @Summary     取得所有聊天室列表
// @Description 取得所有聊天室列表
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                          true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                        true "該 websocket 的操作 [get_all_chat_room]"
// @Param       payload body     sendChatMessageRequestPayload true "type: text"
// @Success     200     {object} message{payload=getAllChatRoomResponse}
// @Failure     400     {object} message
// @Router      /chat/ws/cmd:get_all_chat_room [get]
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

// @Summary     更新特定聊天室的設定 (共通 - 聊天室成員皆套用)
// @Description 更新特定聊天室的設定 (共通 - 聊天室成員皆套用)
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                                true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                              true "該 websocket 的操作 [update_chat_room_setting]"
// @Param       payload body     updateChatRoomSettingRequestPayload true "payload"
// @Success     200     {object} message{payload=updateChatRoomSettingResponse}
// @Failure     400     {object} message
// @Failure     500     {object} message
// @Router      /chat/ws/cmd:update_chat_room_setting [get]
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

// @Summary     更新特定聊天室的設定 (私人 - 僅自己有效)
// @Description 更新特定聊天室的設定 (私人 - 僅自己有效)
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                                      true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                                    true "該 websocket 的操作 [update_chat_room_custom_setting]"
// @Param       payload body     updateChatRoomCustomSettingRequestPayload true "payload"
// @Success     200     {object} message
// @Failure     400     {object} message
// @Failure     500     {object} message
// @Router      /chat/ws/cmd:update_chat_room_custom_setting [get]
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

// @Summary     更新使用者已讀某個聊天室的最後時間
// @Description 更新使用者已讀某個聊天室的最後時間
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                                 true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                               true "該 websocket 的操作 [have_read]"
// @Param       payload body     updateChatRoomLastSeenRequestPayload true "payload"
// @Success     200     {object} message{payload=updateChatRoomCustomSettingResponse}
// @Failure     400     {object} message
// @Failure     500     {object} message
// @Router      /chat/ws/cmd:have_read [get]
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

// @Summary     建立聊天室
// @Description 建立聊天室
// @Tags        chat
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       type   body     string true "聊天室類型 (personal or group)"
// @Param       name   body     string true "Chat Room Name"
// @Param       avatar body     string true "Chat Room Avatar"
// @Success     200    {object} controllers.JSONResponse
// @Failure     400    {object} controllers.JSONResponse
// @Failure     500    {object} controllers.JSONResponse
// @Router      /chat/room [post]
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

	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	chatRoomMemberModel := models.ChatRoomMember{
		ChatRoomID:  chatRoomModel.ID,
		AccountUUID: authPayload.AccountUUID,
		LastSeenAt:  time.Now().UTC(),
	}
	result = chatRoomMemberModel.Create()

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

// @Summary     取得聊天室邀請碼
// @Description 取得聊天室邀請碼
// @Tags        chat
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       chatRoomID path     uint true "Chat Room ID"
// @Success     200        {object} controllers.JSONResponse{data=getRoomInviteCodeResponse}
// @Failure     403        {object} controllers.JSONResponse
// @Failure     500        {object} controllers.JSONResponse
// @Router      /chat/room/:chatRoomID/invite/code [get]
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

// @Summary     藉由邀請碼加入聊天室
// @Description 藉由邀請碼加入聊天室
// @Tags        chat
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       chatRoomID path     uint   true "Chat Room ID"
// @Param       inviteCode path     string true "Chat Room Invite Code"
// @Success     200        {object} controllers.JSONResponse
// @Failure     400        {object} controllers.JSONResponse
// @Failure     500        {object} controllers.JSONResponse
// @Router      /chat/room/:chatRoomID/invite/code/:inviteCode [patch]
func UpdateRoomMemberByInviteCode(c *gin.Context) {
	authPayload := c.MustGet("authPayload").(*models.JWTToken)
	inviteCode := c.Param("inviteCode")

	chatRoomInviteCodeModel := models.ChatRoomInviteCode{
		Code: inviteCode,
	}

	err := chatRoomInviteCodeModel.Read()

	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestChatRoomInviteCodeExpired,
			Message: errMessageRequestChatRoomInviteCodeExpired,
			Data:    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisGetKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	chatRoomMemberModel := &models.ChatRoomMember{
		ChatRoomID:  chatRoomInviteCodeModel.ChatRoomID,
		AccountUUID: authPayload.AccountUUID,
	}
	result := chatRoomMemberModel.Create()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	// TODO: websocket broadcast，發送系統訊息到該聊天室，以此來廣播通知該聊天室的所有成員有新成員的加入
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}
