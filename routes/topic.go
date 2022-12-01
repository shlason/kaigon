package routes

import "github.com/gin-gonic/gin"

func RegisteTopicRoutes(publicR *gin.RouterGroup, privateR *gin.RouterGroup) {
	publicR.GET("/topics")
	privateR.POST("/topics/followed")
}
