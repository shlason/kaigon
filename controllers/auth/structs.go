package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/utils"
)

type getAuthTokenByRefreshTokenRequestParamsPayload struct {
	AccountUUID string `json:"accountUuid"`
	Email       string `json:"email"`
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
