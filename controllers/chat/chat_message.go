package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
)

var acceptCheatMessageTypes = map[string]string{
	"text": "text",
}

type sendChatMessageRequestPayload struct {
	From      string    `json:"from"`
	To        uint      `json:"to"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func (c sendChatMessageRequestPayload) Parse(data interface{}) (sendChatMessageRequestPayload, error) {
	p := sendChatMessageRequestPayload{}

	bytes, err := json.Marshal(data)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(bytes, &p)

	return p, err
}

type sendChatMessageResponsePayload struct {
	From      string    `json:"from"`
	To        uint      `json:"to"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func sendChatMessageHandler(clients map[string]client, msg message) {
	*msg.Self <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds["received"],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}

	sendChatMsgReqPayload, err := sendChatMessageRequestPayload{}.Parse(msg.Payload)

	if err != nil {
		fmt.Println("chatMsgPayload.Parse got error")
		fmt.Println(err)
	}

	chatMsgModel := &models.ChatMessage{
		From:      sendChatMsgReqPayload.From,
		To:        sendChatMsgReqPayload.To,
		Type:      sendChatMsgReqPayload.Type,
		Content:   sendChatMsgReqPayload.Content,
		Timestamp: sendChatMsgReqPayload.Timestamp,
	}

	_, err = chatMsgModel.InsertOne()

	if err != nil {
		fmt.Println("insert one got error: ", err)
	}

	var accountChatRooms *[]models.AccountChatRoom

	result := models.AccountChatRoom{}.ReadAllByChatRoomID(sendChatMsgReqPayload.To, accountChatRooms)

	// TODO: 取得聊天室所有成員資訊時噴錯的處理
	if result.Error != nil {
		fmt.Printf("ReadAllByChatRoomID got error: %s\n", result.Error)
		return
	}

	for _, accountChatRoom := range *accountChatRooms {
		toCli, ok := clients[accountChatRoom.AccountUUID]
		// TODO: 接收方不在線上時的處理
		if !ok {
			fmt.Printf("Friend: %s offline\n", accountChatRoom.AccountUUID)
			return
		}
		fmt.Printf("message sending from: %s, to: %d\n", sendChatMsgReqPayload.From, sendChatMsgReqPayload.To)
		toCli <- message{
			Seq:           msg.Seq,
			Cmd:           acceptResponseCmds["send_chat_message"],
			StatusCode:    http.StatusOK,
			StatusMessage: controllers.SuccessMessage,
			Payload: sendChatMessageResponsePayload{
				From:      sendChatMsgReqPayload.From,
				To:        sendChatMsgReqPayload.To,
				Content:   sendChatMsgReqPayload.Content,
				Timestamp: time.Now().UTC(),
			},
		}
		fmt.Printf("message sended from: %s, to: %d\n", sendChatMsgReqPayload.From, sendChatMsgReqPayload.To)
	}
}
