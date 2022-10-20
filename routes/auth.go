package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/auth"
)

func RegisteAuthRoutes(publicR *gin.RouterGroup) {
	// OAuth callback
	publicR.GET("/auth/o/google/url", auth.OAuthCallbackForGoogle)
	publicR.GET("/auth/o/google/login", auth.OAuthCallbackForGoogle)
	publicR.GET("/auth/o/google/bind", auth.OAuthCallbackForGoogle)

	// Session
	publicR.GET("/auth/session/token/refresh", auth.GetAuthTokenByRefreshToken)

	// Captcha
	publicR.GET("/auth/captcha", auth.GetCaptchaInfo)
	publicR.GET("/auth/captcha/:captchaUUID/image", auth.GetCaptchaImage)
	publicR.GET("/auth/captcha/:captchaUUID/refresh", auth.UpdateCaptchaInfo)
}
