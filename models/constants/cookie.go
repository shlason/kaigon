package constants

type refreshTokenCookieInfo struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

var RefreshTokenCookieInfo = refreshTokenCookieInfo{
	Name: "REFRESH_TOKEN",
	// Second Base -> 分 -> 小時 -> 天 -> 20 天
	MaxAge:   60 * 60 * 24 * 20,
	Path:     "/",
	Domain:   "/",
	Secure:   false,
	HttpOnly: true,
}
