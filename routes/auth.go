package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/auth"
)

func RegisteAuthRoutes(publicR *gin.RouterGroup) {
	// OAuth callback
	publicR.GET("/auth/o/google/callback", auth.OAuthCallbackForGoogle)

	// Captcha
	publicR.GET("/auth/captcha", auth.GetCaptchaInfo)
	publicR.GET("/auth/captcha/:captchaUUID/image", auth.GetCaptchaImage)
	publicR.GET("/auth/captcha/:captchaUUID/refresh", auth.UpdateCaptchaInfo)
}
