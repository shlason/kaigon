package account

import (
	"fmt"

	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
)

const passwordMaximumLength int = 32

type signUpRequestPayload struct {
	Email       string
	Password    string
	CaptchaUUID string
	CaptchaCode string
}

func (p *signUpRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
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

	if p.CaptchaUUID == "" || p.CaptchaCode == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadCaptchaFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadCaptchaFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type signInRequestPayload struct {
	Email       string
	Password    string
	CaptchaUUID string
	CaptchaCode string
}

func (p *signInRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
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

	if p.CaptchaUUID == "" || p.CaptchaCode == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadCaptchaFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadCaptchaFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type signInResponsePayload struct {
	Token string
}

type createResetPasswordSessionRequestPayload struct {
	Email string
}

func (p *createResetPasswordSessionRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	if !utils.IsValidEmailAddress(p.Email) {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadEmailFieldNotValid,
			Message: errMessageRequestPayloadEmailFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type resetPasswordTemplateParams struct {
	Link string
}

func (p *resetPasswordTemplateParams) generate(accountUUID string) error {
	authAccountResetPasswordModel := &models.AuthAccountResetPassword{
		AccountUUID: accountUUID,
	}
	err := authAccountResetPasswordModel.Create()
	if err != nil {
		return err
	}
	p.Link = fmt.Sprintf("{{ResetPasswordPageURL}}?token=%s&code=%s", authAccountResetPasswordModel.Token, authAccountResetPasswordModel.Code)

	return err
}

type resetPasswordRequestPayload struct {
	Email       string
	Password    string
	Token       string
	Code        string
	CaptchaUUID string
	CaptchaCode string
}

func (p *resetPasswordRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
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

	if p.CaptchaUUID == "" || p.CaptchaCode == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadCaptchaFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadCaptchaFieldNotValid,
			Data:    nil,
		}, true
	}

	if p.Token == "" || p.Code == "" {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadTokenCodeFieldsNotValid,
			Message: errMessageRequestPayloadTokenCodeFieldsNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}
