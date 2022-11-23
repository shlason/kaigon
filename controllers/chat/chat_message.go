package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func sendChatMessageHandler(clients map[string]client, msg message) {
	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds[acceptResponseCmds["received"]],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}

	sendChatMsgReqPayload, err := sendChatMessageRequestPayload{}.parse(msg.Payload)

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

func getChatMessage(msg message) {
	getChatMsgReqPayload, err := getChatMessageRequestPayload{}.Parse(msg.Payload)

	if err != nil {
		fmt.Println("getChatMessageRequestPayload.Parse got error")
		fmt.Println(err)
		return
	}

	chatMsgsResp := chatMessagesResponse{}

	// TODO: pagination
	chateMessages, err := models.ChatMessage{}.FindByTo(getChatMsgReqPayload.ChatRoomID)

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
