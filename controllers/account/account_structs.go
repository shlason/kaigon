package account

import (
	"fmt"

	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
)

const passwordMaximumLength int = 32

type signUpRequestPayload struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CaptchaUUID string `json:"captchaUuid"`
	CaptchaCode string `json:"captchaCode"`
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
	Email       string `json:"email"`
	Password    string `json:"password"`
	CaptchaUUID string `json:"captchaUuid"`
	CaptchaCode string `json:"captchaCode"`
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

type createResetPasswordSessionRequestPayload struct {
	Email string `json:"email"`
	Path  string `json:"path"`
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
	Link string `json:"link"`
}

func (p *resetPasswordTemplateParams) generate(accountUUID, email, path string) error {
	authAccountResetPasswordModel := &models.AuthAccountResetPassword{
		AccountUUID: accountUUID,
	}
	err := authAccountResetPasswordModel.Create()
	if err != nil {
		return err
	}
	p.Link = fmt.Sprintf(
		"%s://%s%s?email=%s&token=%s&code=%s",
		configs.Server.Protocol,
		configs.Server.Host,
		path,
		email,
		authAccountResetPasswordModel.Token,
		authAccountResetPasswordModel.Code,
	)

	return nil
}

type resetPasswordRequestPayload struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	Code        string `json:"code"`
	CaptchaUUID string `json:"captchaUuid"`
	CaptchaCode string `json:"captchaCode"`
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

type verificationSessionTemplateParams struct {
	Link string `json:"link"`
}

func (p *verificationSessionTemplateParams) generate(accountUUID string) error {
	authAccountEmailVerificationModel := &models.AuthAccountEmailVerification{
		AccountUUID: accountUUID,
	}
	err := authAccountEmailVerificationModel.Create()
	if err != nil {
		return err
	}
	p.Link = fmt.Sprintf(
		"%s://%s/account/%s/info/verification/email?token=%s&code=%s",
		configs.Server.Protocol,
		configs.Server.Host,
		accountUUID,
		authAccountEmailVerificationModel.Token,
		authAccountEmailVerificationModel.Code,
	)

	return nil
}

type createVerifySessionRequestPayload struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

func (p *createVerifySessionRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	var acceptVerificationType = map[string]string{
		"email": "email",
	}

	if !utils.IsValidEmailAddress(p.Email) {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadEmailFieldNotValid,
			Message: errMessageRequestPayloadEmailFieldNotValid,
			Data:    nil,
		}, true
	}

	if _, ok := acceptVerificationType[p.Type]; !ok {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadVerificationTypeFieldNotValid,
			Message: errMessageRequestPayloadVerificationTypeFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type verifyWithEmailRequestPayload struct {
	AccountUUID string
	Token       string
	Code        string
}
