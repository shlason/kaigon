package chat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shlason/kaigon/controllers"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary     建立 Chat Websocket 連線 (HTTP)
// @Description HTTP GET 方法，在處理請求時會切換協議由 HTTP -> Websocket
// @Tags        chat
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} controllers.JSONResponse
// @Failure     500 {object} controllers.JSONResponse
// @Router      /chat/ws [get]
func Connect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	// TODO: 屆時改為 authPayload 來源
	accountUUID := c.Query("accountUuid")
	cli := make(client)

	connInfo := connectionInfo{
		client:      &cli,
		AccountUUID: accountUUID,
	}
	clientConnect <- connInfo

	defer (func() {
		ws.Close()
		clientDisconnect <- connInfo
	})()

	go func() {
		for msg := range cli {
			fmt.Printf("Actually get client: %s channel message: %v\n", connInfo.AccountUUID, msg)

			fmt.Printf("Client: %s writing JSON to websocket\n", connInfo.AccountUUID)
			err := ws.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Client: %s write json got error: %s\n", connInfo.AccountUUID, err)
				break
			}
			fmt.Printf("Client: %s writed JSON to websocket\n", connInfo.AccountUUID)

			fmt.Printf("Waiting for client: %s channel to get message\n", connInfo.AccountUUID)
		}
	}()

	for {
		var msg message

		fmt.Printf("Client: %s waiting JSON from websocket\n", connInfo.AccountUUID)

		err := ws.ReadJSON(&msg)

		fmt.Printf("Client: %s reading JSON from websocket\n", connInfo.AccountUUID)

		if err != nil {
			break
		}

		msg.Self = selfInfo{
			Channel:     &cli,
			AccountUUID: accountUUID,
		}

		fmt.Printf("Client: %s passing msg to messages channel\n", connInfo.AccountUUID)
		messages <- msg
		fmt.Printf("Client: %s passed msg to messages channel\n", connInfo.AccountUUID)
	}
}
