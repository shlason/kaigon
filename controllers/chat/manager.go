package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shlason/kaigon/controllers"
)

type client chan message

type connectionInfo struct {
	*client
	AccountUUID string
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

				chatMsgPayload, err := chatMessagePayload{}.Parse(msg.Payload)

				if err != nil {
					fmt.Println("chatMsgPayload.Parse got error")
					fmt.Println(err)
				}
				fmt.Println(chatMsgPayload)
				toCli, ok := clients[chatMsgPayload.To]
				// TODO: 接收方不在線上時的處理
				if !ok {
					fmt.Printf("Friend: %s offline\n", chatMsgPayload.To)
					continue
				}
				fmt.Printf("message sending from: %s, to: %s\n", chatMsgPayload.From, chatMsgPayload.To)
				toCli <- message{
					Seq:           msg.Seq,
					Cmd:           acceptCmds["chat_message"],
					StatusCode:    http.StatusOK,
					StatusMessage: controllers.SuccessMessage,
					Payload: chatMessagePayload{
						From:      chatMsgPayload.From,
						To:        chatMsgPayload.To,
						Text:      chatMsgPayload.Text,
						Timestamp: time.Now().UTC(),
					},
				}
				fmt.Printf("message sended from: %s, to: %s\n", chatMsgPayload.From, chatMsgPayload.To)
			}

		case connInfo := <-clientConnect:
			fmt.Printf("%s connecting\n", connInfo.AccountUUID)
			clients[connInfo.AccountUUID] = *connInfo.client
		case connInfo := <-clientDisconnect:
			fmt.Printf("%s disconnect\n", connInfo.AccountUUID)
			delete(clients, connInfo.AccountUUID)
			close(*connInfo.client)
		}
	}
}
