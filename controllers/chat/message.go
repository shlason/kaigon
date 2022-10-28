package chat

var acceptCmds = map[string]string{
	"ping":         "ping",
	"pong":         "pong",
	"chat_message": "chat_message",
	"received":     "received",
}

// TODO: check validation
type message struct {
	Self          *client `json:",omitempty"`
	Seq           int
	Cmd           string `json:"cmd"`
	StatusCode    int    `json:"statusCode,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Payload       interface{}
}
