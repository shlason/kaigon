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

	publicR.GET("/account/:accountUUID/info/verification/email", account.VerifyWithEmail)

	privateR.POST("/account/:accountUUID/info/verification", account.CreateVerifySession)
}
