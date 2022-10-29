package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/chat"
)

func RegisteChatRoutes(r *gin.RouterGroup) {
	r.GET("/chat", chat.Connect)
}
