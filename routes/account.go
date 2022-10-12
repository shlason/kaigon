package routes

import "github.com/gin-gonic/gin"

func RegisteAccountRoutes(r *gin.RouterGroup) {
	r.POST("/account/signup")
	r.POST("/account/signin")
	r.POST("/account/:accountID/verification")
	r.GET("/account/:accountID/verification/email")
	r.POST("/account/:accountID/settings/password/reset")
	r.PATCH("/account/:accountID/settings/password/reset")
}
