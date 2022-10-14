package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/account"
)

func RegisteAccountRoutes(r *gin.RouterGroup) {
	r.POST("/account/signup", account.SignUp)
	r.POST("/account/signin", account.SignIn)
	r.POST("/account/info/password/reset", account.CreateResetPasswordSession)
	r.PATCH("/account/info/password/reset", account.ResetPassword)
	r.POST("/account/:accountID/info/verification", account.CreateVerifySession)
	r.GET("/account/:accountID/info/verification/email", account.VerifyWithEmail)
}
