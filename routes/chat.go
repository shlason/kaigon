package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/chat"
)

func RegisteChatRoutes(r *gin.RouterGroup) {
	r.GET("/chat/ws/:token", chat.Connect)
	// Create chat room
	r.POST("/chat/room", chat.CreateRoom)
	// Get chat room invite code
	r.GET("/chat/room/:chatRoomID/invite/code", chat.GetRoomInviteCode)
	// Join chat room by invite code
	r.PATCH("/chat/room/:chatRoomID/invite/code/:inviteCode", chat.UpdateRoomMemberByInviteCode)
}
