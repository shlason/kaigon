package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/auth"
)

func RegisteAuthRoutes(r *gin.RouterGroup) {
	// OAuth callback
	r.GET("/auth/o/google/callback", auth.OAuthCallbackForGoogle)

	// Captcha
	r.GET("/auth/captcha", auth.GetCaptchaInfo)
	r.GET("/auth/captcha/:captchaUUID/image", auth.GetCaptchaImage)
	r.GET("/auth/captcha/:captchaUUID/refresh", auth.UpdateCaptchaInfo)
}
