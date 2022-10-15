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

	privateR.POST("/account/:accountID/info/verification", account.CreateVerifySession)
	privateR.GET("/account/:accountID/info/verification/email", account.VerifyWithEmail)
}
