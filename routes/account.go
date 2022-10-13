package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/account"
)

func RegisteAccountRoutes(r *gin.RouterGroup) {
	r.POST("/account/signup", account.SignUp)
	r.POST("/account/signin", account.SignIn)
	r.POST("/account/:accountID/verification", account.CreateVerifySession)
	r.GET("/account/:accountID/verification/email", account.VerifyWithEmail)
	r.POST("/account/:accountID/settings/password/reset", account.CreateResetPasswordSession)
	r.PATCH("/account/:accountID/settings/password/reset", account.ResetPassword)
}
