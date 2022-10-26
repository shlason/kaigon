package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/account"
)

func RegisteAccountRoutes(publicR *gin.RouterGroup, privateR *gin.RouterGroup) {
	publicR.POST("/account/signup", account.SignUp)
	publicR.POST("/account/signin", account.SignIn)

	publicR.POST("/account/info/password/reset", account.CreateResetPasswordSession)
	publicR.PATCH("/account/info/password/reset", account.ResetPassword)

	privateR.POST("/account/:accountUUID/info/verification", account.CreateVerifySession)
	publicR.GET("/account/:accountUUID/info/verification/email", account.VerifyWithEmail)

	publicR.GET("/account/:accountUUID/profile", account.GetProfile)
	privateR.PATCH("/account/:accountUUID/profile")

	publicR.GET("/account/:accountUUID/setting")
	privateR.PATCH("/account/:accountUUID/setting")

	privateR.GET("/account/:accountUUID/setting/notification")
	privateR.PUT("/account/:accountUUID/setting/notification")
}
