package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/utils"
)

type getOAuthUrlQueryParmas struct {
	Type string `form:"type"`
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

type getAuthTokenByRefreshTokenRequestParamsPayload struct {
	AccountUUID string `form:"accountUuid"`
	Email       string `form:"email"`
}

func (p *getAuthTokenByRefreshTokenRequestParamsPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	if utils.IsValidEmailAddress(p.Email) {
		return controllers.JSONResponse{}, false
	}
	return controllers.JSONResponse{
		Code:    ErrCodeRequestQueryParamEmailFieldNotValid,
		Message: ErrMessageRequestQueryParamEmailFieldNotValid,
		Data:    nil,
	}, true
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
