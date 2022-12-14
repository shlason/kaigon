package chat

import (
	"fmt"
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

	// TODO: 增加針對聊天室全體成員的廣播
	for {
		select {
		case msg := <-messages:
			// TODO: 回傳的 status message, code 需要討論和統整
			// TODO: message 的 interface{} type payload 每次 parse 成相對應的 struct type 都需要重複寫 parse method 因此要想辦法抽出來複用
			errResp, isNotValid := msg.Check()

			if isNotValid {
				*msg.Self.Channel <- errResp
				return
			}

			switch msg.Cmd {
			case acceptRequestCmds["ping"]:
				PingHandler(msg)
			case acceptRequestCmds["get_all_chat_room"]:
				getAllChatRoomHandler(msg)
			case acceptRequestCmds["get_chat_message"]:
				getChatMessage(msg)
			case acceptRequestCmds["send_chat_message"]:
				sendChatMessageHandler(clients, msg)
			case acceptRequestCmds["update_chat_room_setting"]:
				updateChatRoomSettingHandler(clients, msg)
			case acceptRequestCmds["update_chat_room_custom_setting"]:
				updateChatRoomCustomSettingHandler(msg)
			case acceptRequestCmds["have_read"]:
				updateChatRoomLastSeenHandler(clients, msg)
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
