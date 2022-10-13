package routes

import "github.com/gin-gonic/gin"

func RegisteAuthRoutes(r *gin.RouterGroup) {
	r.GET("/auth/o/google/callback")

	r.GET("/auth/captcha")
	r.GET("/auth/captcha/:captchaUUID/image")
	r.GET("/auth/captcha/:captchaUUID/refresh")
}
