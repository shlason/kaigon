package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/auth"
)

func RegisteAuthRoutes(publicR *gin.RouterGroup, privateR *gin.RouterGroup) {
	// OAuth callback
	publicR.GET("/auth/o/google/url", auth.GetGoogleOAuthURL)
	publicR.GET("/auth/o/google/login", auth.GoogleOAuthRedirectURIForLogin)
	publicR.GET("/auth/o/google/bind", auth.GoogleOAuthRedirectURIForBind)
	privateR.PATCH("/auth/o/google/bind", auth.GoogleOAuthBind)

	// Session
	publicR.GET("/auth/session/token/refresh", auth.GetAuthTokenByRefreshToken)

	// Captcha
	publicR.GET("/auth/captcha", auth.GetCaptchaInfo)
	publicR.GET("/auth/captcha/:captchaUUID/image", auth.GetCaptchaImage)
	publicR.GET("/auth/captcha/:captchaUUID/refresh", auth.UpdateCaptchaInfo)
}
