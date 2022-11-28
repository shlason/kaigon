package chat

import (
	"net/http"

	"github.com/shlason/kaigon/controllers"
)

// @Summary     心跳包
// @Description 心跳包
// @Tags        chat - websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       seq body     uint   true "使用 timestamp 以此表示個別 message 的辨識 id"
// @Param       cmd body     string true "該 websocket 的操作 [ping]"
// @Success     200 {object} message
// @Router      /chat/ws/cmd:ping [get]
func PingHandler(msg message) {
	*msg.Self.Channel <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds["pong"],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}
}
