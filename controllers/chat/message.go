package chat

import "net/http"

var acceptRequestCmds = map[string]string{
	"ping":                             "ping",
	"get_all_chat_room":                "get_all_chat_room",
	"get_chat_message":                 "get_chat_message",
	"send_chat_message":                "send_chat_message",
	"get_chat_room_setting":            "get_chat_room_setting",
	"update_chat_room_setting":         "update_chat_room_setting",
	"update_chat_room_account_setting": "update_chat_room_account_setting",
}

var acceptResponseCmds = map[string]string{
	"pong":                             "pong",
	"get_all_chat_room":                "get_all_chat_room",
	"get_chat_message":                 "get_chat_message",
	"send_chat_message":                "send_chat_message",
	"get_chat_room_setting":            "get_chat_room_setting",
	"update_chat_room_setting":         "update_chat_room_setting",
	"update_chat_room_account_setting": "update_chat_room_account_setting",
	"received":                         "received",
}

type message struct {
	Self          *client     `json:",omitempty"`
	Seq           int         `json:"seq"`
	Cmd           string      `json:"cmd"`
	CustomCode    string      `json:"customCode"`
	StatusCode    int         `json:"statusCode"`
	StatusMessage string      `json:"statusMessage"`
	Payload       interface{} `json:"payload"`
}

func (m message) Check() (errResponseMsg message, isNotValid bool) {
	if _, ok := acceptRequestCmds[m.Cmd]; !ok {
		return message{
			Seq:           m.Seq,
			Cmd:           m.Cmd,
			CustomCode:    errCodeRequestFieldNotValid,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: errMessageRequestFieldNotValid,
			Payload:       nil,
		}, true
	}
	return message{}, false
}
