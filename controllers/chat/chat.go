package chat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Connect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// TODO: 屆時改為 authPayload 來源
	accountUUID := c.Query("accountUuid")
	fmt.Println(accountUUID)
	cli := make(client)

	connInfo := connectionInfo{
		AccountUUID: accountUUID,
		client:      &cli,
	}
	clientConnect <- connInfo

	defer (func() {
		ws.Close()
		clientDisconnect <- connInfo
	})()

	for {
		var msg message

		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}

		msg.Self = &cli
		messages <- msg

		msg = <-cli

		ws.WriteJSON(msg)

		if err != nil {
			break
		}
	}
}
