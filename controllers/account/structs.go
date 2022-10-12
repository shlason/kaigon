package account

import (
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/utils"
)

type signUpRequestPayload struct {
	Email       string
	Password    string
	CaptchaUUID string
	CaptchaVal  string
}

func (p *signUpRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	const passwordMaximumLength int = 32

	if !utils.IsValidEmailAddress(p.Email) {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadEmailFieldNotValid,
			Message: errMessageRequestPayloadEmailFieldNotValid,
			Data:    nil,
		}, true
	}
	if p.Password == "" || len(p.Password) > passwordMaximumLength {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadPasswordFieldNotValid,
			Message: errMessageRequestPayloadPasswordFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type signInRequestPayload struct {
	Email       string
	Password    string
	CaptchaUUID string
	CaptchaVal  string
}
