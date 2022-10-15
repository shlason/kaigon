package routes

import "github.com/gin-gonic/gin"

func RegisteAuthRoutes(publicR *gin.RouterGroup) {
	publicR.GET("/auth/o/google/callback")

	publicR.GET("/auth/captcha")
	publicR.GET("/auth/captcha/:captchaUUID/image")
	publicR.GET("/auth/captcha/:captchaUUID/refresh")
}
