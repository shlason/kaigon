package chat

import (
	"fmt"
	"net/http"

	"github.com/shlason/kaigon/controllers"
)

type client chan message

type connectionInfo struct {
	AccountUUID string
	*client
}

var (
	clientConnect    = make(chan connectionInfo)
	clientDisconnect = make(chan connectionInfo)
	messages         = make(chan message)
)

func init() {
	go clientManager()
}

func clientManager() {
	// map[AccountUUID]client
	clients := make(map[string]client)

	for {
		select {
		case msg := <-messages:
			fmt.Println(msg)
			// TODO: 回傳的 status message, code 需要討論和統整
			if msg.Cmd == acceptCmds["ping"] {
				*msg.Self <- message{
					Seq:           msg.Seq,
					Cmd:           acceptCmds["pong"],
					StatusCode:    http.StatusOK,
					StatusMessage: controllers.SuccessMessage,
					Payload:       nil,
				}
				continue
			} else if msg.Cmd == acceptCmds["chat_message"] {
				*msg.Self <- message{
					Seq:           msg.Seq,
					Cmd:           acceptCmds["received"],
					StatusCode:    http.StatusOK,
					StatusMessage: controllers.SuccessMessage,
					Payload:       nil,
				}

				var anyPayload interface{} = msg.Payload
				var chatMsgPayload chatMessagePayload

				payload := anyPayload.(*chatMessagePayload)
				chatMsgPayload.From = payload.From
				chatMsgPayload.To = payload.To
				chatMsgPayload.Text = payload.Text
				chatMsgPayload.Timestamp = payload.Timestamp

				toCli, ok := clients[chatMsgPayload.To]
				// TODO: 接收方不在線上時的處理
				if !ok {
					fmt.Println("Friend offline")
					continue
				}
				toCli <- message{
					Seq:           msg.Seq,
					Cmd:           acceptCmds["chat_message"],
					StatusCode:    http.StatusOK,
					StatusMessage: controllers.SuccessMessage,
					Payload:       chatMsgPayload,
				}

				continue
			}

		case connInfo := <-clientConnect:
			clients[connInfo.AccountUUID] = *connInfo.client
		case connInfo := <-clientDisconnect:
			delete(clients, connInfo.AccountUUID)
			close(*connInfo.client)
		}
	}
}
