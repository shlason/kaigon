package routes

import "github.com/gin-gonic/gin"

func RegisteSearchRoutes(publicR *gin.RouterGroup) {
	publicR.GET("/search/forums")
	publicR.GET("/search/posts")
	publicR.GET("/search/topics")
	publicR.GET("/search/accounts")
}
