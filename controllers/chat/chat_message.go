package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary     傳送聊天訊息
// @Description 傳送聊天訊息
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                          true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                        true "該 websocket 的操作 [send_chat_message]"
// @Param       payload body     sendChatMessageRequestPayload true "type: text"
// @Success     200     {object} message{payload=chatMessageResponse}
// @Failure     400     {object} message
// @Router      /chat/ws/cmd:send_chat_message [get]
func sendChatMessageHandler(clients map[string]client, msg message) {
	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds[acceptResponseCmds["received"]],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}

	sendChatMsgReqPayload, err := sendChatMessageRequestPayload{}.parse(msg.Payload)

	// TODO: error handle
	if err != nil {
		fmt.Println("chatMsgPayload.Parse got error")
		fmt.Println(err)
		return
	}

	isNotValid := sendChatMsgReqPayload.check()

	if isNotValid {
		*msg.Self.Channel <- message{
			Seq:           msg.Seq,
			Cmd:           msg.Cmd,
			CustomCode:    errCodeRequestFieldNotValid,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: errMessageRequestFieldNotValid,
			Payload:       nil,
		}
		return
	}

	chatMsgModel := &models.ChatMessage{
		From:      sendChatMsgReqPayload.From,
		To:        sendChatMsgReqPayload.To,
		Type:      sendChatMsgReqPayload.Type,
		Content:   sendChatMsgReqPayload.Content,
		Timestamp: time.Now().UTC(),
	}

	mgResult, err := chatMsgModel.InsertOne()

	if err != nil {
		fmt.Println("insert one got error: ", err)
		return
	}

	var chatRoomMembers []models.ChatRoomMember

	result := models.ChatRoomMember{}.ReadAllByChatRoomID(sendChatMsgReqPayload.To, &chatRoomMembers)

	// TODO: 取得聊天室所有成員資訊時噴錯的處理
	if result.Error != nil {
		fmt.Printf("ReadAllByChatRoomID got error: %s\n", result.Error)
		return
	}
	fmt.Println(chatRoomMembers)
	for _, chatRoomMember := range chatRoomMembers {
		toCli, ok := clients[chatRoomMember.AccountUUID]
		// TODO: 接收方不在線上時的處理
		if !ok {
			fmt.Printf("Friend: %s offline\n", chatRoomMember.AccountUUID)
			continue
		}
		fmt.Printf("message sending from: %s, to: %d\n", sendChatMsgReqPayload.From, sendChatMsgReqPayload.To)
		toCli <- message{
			Seq:           msg.Seq,
			Cmd:           acceptResponseCmds[acceptResponseCmds["send_chat_message"]],
			StatusCode:    http.StatusOK,
			StatusMessage: controllers.SuccessMessage,
			Payload: chatMessageResponse{
				ID:        mgResult.InsertedID.(primitive.ObjectID).Hex(),
				From:      sendChatMsgReqPayload.From,
				To:        sendChatMsgReqPayload.To,
				Type:      sendChatMsgReqPayload.Type,
				Content:   sendChatMsgReqPayload.Content,
				Timestamp: time.Now().UTC(),
			},
		}
		fmt.Printf("message sended from: %s, to: %d\n", sendChatMsgReqPayload.From, sendChatMsgReqPayload.To)
	}
}

// @Summary     取得某聊天室中的所有訊息
// @Description 取得某聊天室中的所有訊息
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq     body     uint                         true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd     body     string                       true "該 websocket 的操作 [get_chat_message]"
// @Param       payload body     getChatMessageRequestPayload true "Chat Room ID"
// @Success     200     {object} message{payload=chatMessagesResponse}
// @Failure     400     {object} message
// @Router      /chat/ws/cmd:get_chat_message [get]
func getChatMessage(msg message) {
	getChatMsgReqPayload, err := getChatMessageRequestPayload{}.Parse(msg.Payload)

	// TODO: error handle
	if err != nil {
		fmt.Println("getChatMessageRequestPayload.Parse got error")
		fmt.Println(err)
		return
	}

	chatMsgsResp := chatMessagesResponse{}

	// TODO: pagination
	chateMessages, err := models.ChatMessage{}.FindByTo(getChatMsgReqPayload.ChatRoomID)

	// TODO: error handle
	if err != nil {
		fmt.Println("ChatMessage{}.FindByTo got error")
		fmt.Println(err)
		return
	}

	for _, chatMsg := range chateMessages {
		chatMsgsResp = append(chatMsgsResp, chatMessageResponse{
			ID:        chatMsg.ID.Hex(),
			From:      chatMsg.From,
			To:        chatMsg.To,
			Type:      chatMsg.Type,
			Content:   chatMsg.Content,
			Timestamp: chatMsg.Timestamp,
		})
	}

	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds[acceptResponseCmds["get_chat_message"]],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       chatMsgsResp,
	}
}
