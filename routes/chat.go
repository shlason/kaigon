package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/chat"
)

func RegisteChatRoutes(r *gin.RouterGroup) {
	r.GET("/chat/ws", chat.Connect)
	// Create chat room
	r.POST("/chat/room")
	// Get chat room invite code
	r.GET("/chat/room/:chatRoomID/invite/code")
	// Join chat room by invite code
	r.PATCH("/chat/room/:chatRoomID/invite/code/:inviteCode")
	// Update chat room related setting
	r.PATCH("/chat/room/:chatRoomID/setting")
	r.PATCH("/chat/room/:chatRoomID/account/setting")
}
