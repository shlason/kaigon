package chat

import (
	"net/http"

	"github.com/shlason/kaigon/controllers"
)

func PingHandler(msg message) {
	*msg.Self <- message{
		Seq:           msg.Seq,
		Cmd:           acceptResponseCmds["pong"],
		StatusCode:    http.StatusOK,
		StatusMessage: controllers.SuccessMessage,
		Payload:       nil,
	}
}
