package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	// TODO: 限制規則 regex
	config.AllowOrigins = []string{"https://google.com", "http://local.kaigon.sidesideeffect.io:8080", "http://local.kaigon.sidesideeffect.io:3000"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS", "PUT"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	config.AllowCredentials = true

	return cors.New(config)
}
