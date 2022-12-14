package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
)

type getOAuthUrlQueryParmas struct {
	Type         string `form:"type"`
	RedirectPath string `form:"redirectPath"`
}

func (p *getOAuthUrlQueryParmas) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	var acceptOAuthUrlType = map[string]string{
		"login": "login",
		"bind":  "bind",
	}

	if _, ok := acceptOAuthUrlType[p.Type]; !ok {
		return controllers.JSONResponse{
			Code:    ErrCodeRequestQueryParamsOAuthURLTypeFieldNotValid,
			Message: ErrMessageRequestQueryParamsOAuthURLTypeFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type getOAuthUrlResponsePayload struct {
	URL string `json:"url"`
}

type googleOAuthAccessTokenResponsePayload struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type googleOAuthUserInfoResponsePayload struct {
	Email string `json:"email"`
}

type googleOAuthRedirectURIForLoginQueryParmas struct {
	Code  string `form:"code"`
	State string `form:"form"`
}

type googleOAuthRedirectURIForBindQueryParams struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

type googleOAuthBindRequestPayload struct {
	GrantCode string `json:"grantCode"`
}

type getAuthTokenByRefreshTokenResponsePayload struct {
	AuthToken string `json:"authToken"`
}

type getCaptchaInfoResponsePayload struct {
	UUID string `json:"uuid"`
}

type getCaptchaImageRequestParamsPayload struct {
	CaptchaUUID string `json:"captchaUuid"`
}

func (p *getCaptchaImageRequestParamsPayload) BindParams(c *gin.Context) {
	p.CaptchaUUID = c.Param("captchaUUID")
}

type updateCaptchaImageRequestParamsPayload struct {
	CaptchaUUID string `json:"captchaUuid"`
}

func (p *updateCaptchaImageRequestParamsPayload) BindParams(c *gin.Context) {
	p.CaptchaUUID = c.Param("captchaUUID")
}

type getChatWSTokenResponse struct {
	Token string `json:"token"`
}
