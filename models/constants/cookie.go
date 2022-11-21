package constants

import (
	"fmt"

	"github.com/shlason/kaigon/configs"
)

type refreshTokenCookieInfo struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite int
}

var RefreshTokenCookieInfo = refreshTokenCookieInfo{
	Name: "REFRESH_TOKEN",
	// Second Base -> 分 -> 小時 -> 天 -> 20 天
	MaxAge:   60 * 60 * 24 * 20,
	Path:     "/",
	Domain:   fmt.Sprintf(".%s", configs.Server.Host),
	Secure:   true,
	HttpOnly: true,
	SameSite: 1,
}
