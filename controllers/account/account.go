package account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	var requestPayload *signUpRequestPayload
	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(requestPayload.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	authCaptchaModel := &models.AuthCaptcha{
		UUID: requestPayload.CaptchaUUID,
	}
	err = authCaptchaModel.ReadByUUID()
	if err != nil || authCaptchaModel.Code != requestPayload.CaptchaCode {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadCaptchaFieldCompareMismatch,
			Message: errMessageRequestPayloadCaptchaFieldCompareMismatch,
			Data:    nil,
		})
		return
	}
	authCaptchaModel.Delete()
	accountModel := &models.Account{
		UUID:     uuid.NewString(),
		Email:    requestPayload.Email,
		Password: string(hashPwd),
	}
	result := accountModel.ReadByEmail()
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, controllers.JSONResponse{
			Code:    errCodeRequestPayloadEmailFieldDatabaseRecordAlreadyExist,
			Message: errMessageRequestPayloadEmailFieldDatabaseRecordAlreadyExist,
			Data:    nil,
		})
		return
	}
	result = accountModel.Create()
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

func SignIn(c *gin.Context) {
	var requestPayload *signInRequestPayload
	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	authCaptchaModel := &models.AuthCaptcha{
		UUID: requestPayload.CaptchaUUID,
	}
	err = authCaptchaModel.ReadByUUID()
	if err != nil || authCaptchaModel.Code != requestPayload.CaptchaCode {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadCaptchaFieldCompareMismatch,
			Message: errMessageRequestPayloadCaptchaFieldCompareMismatch,
			Data:    nil,
		})
		return
	}
	authCaptchaModel.Delete()
	accountModel := &models.Account{
		Email: requestPayload.Email,
	}
	result := accountModel.ReadByEmail()
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, controllers.JSONResponse{
				Code:    errCodeRequestPayloadEmailFieldDatabaseRecordNotFound,
				Message: errMessageRequestPayloadEmailFieldDatabaseRecordNotFound,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(accountModel.Password), []byte(requestPayload.Password)) != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadPasswordFieldCompareMismatch,
			Message: errMessageRequestPayloadPasswordFieldCompareMismatch,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

func CreateVerifySession(c *gin.Context) {

}

func VerifyWithEmail(c *gin.Context) {

}

func CreateResetPasswordSession(c *gin.Context) {
	var requestPayload *createResetPasswordSessionRequestPayload
	errResponse, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	fmt.Println(requestPayload)
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	accountModel := &models.Account{
		Email: requestPayload.Email,
	}
	result := accountModel.ReadByEmail()
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, controllers.JSONResponse{
				Code:    errCodeRequestPayloadEmailFieldDatabaseRecordNotFound,
				Message: errMessageRequestPayloadEmailFieldDatabaseRecordNotFound,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}
	to := []string{
		requestPayload.Email,
	}
	templatParams := &resetPasswordTemplateParams{}
	err = templatParams.generate(accountModel.UUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisSetNXKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	err = utils.SendEmail(to, "[Kaigon]：變更 Kaigon 帳號密碼的操作指示", "reset_password.html", templatParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerSendEmailGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

func ResetPassword(c *gin.Context) {
	var requestPayload *resetPasswordRequestPayload
	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	// Captcha auth
	authCaptchaModel := &models.AuthCaptcha{
		UUID: requestPayload.CaptchaUUID,
	}
	err = authCaptchaModel.ReadByUUID()
	if err != nil || authCaptchaModel.Code != requestPayload.CaptchaCode {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadCaptchaFieldCompareMismatch,
			Message: errMessageRequestPayloadCaptchaFieldCompareMismatch,
			Data:    nil,
		})
		return
	}
	authCaptchaModel.Delete()

	// Account check
	accountModel := &models.Account{Email: requestPayload.Email}
	result := accountModel.ReadByEmail()
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, controllers.JSONResponse{
				Code:    errCodeRequestPayloadEmailFieldDatabaseRecordNotFound,
				Message: errMessageRequestPayloadEmailFieldDatabaseRecordNotFound,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	// Token, Code auth
	authAccountResetPasswordModel := &models.AuthAccountResetPassword{
		AccountUUID: accountModel.UUID,
		Token:       requestPayload.Token,
		Code:        requestPayload.Code,
	}
	err = authAccountResetPasswordModel.Read()
	fmt.Println("Resid Val", authAccountResetPasswordModel.Result)
	if !(err == nil && authAccountResetPasswordModel.IsMatch()) {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadTokenCodeFieldsNotValid,
			Message: errMessageRequestPayloadTokenCodeFieldsNotValid,
			Data:    nil,
		})
		return
	}

	// Modify account password
	result = accountModel.UpdatePasswordByEmail(requestPayload.Password)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseUpdateGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	authAccountResetPasswordModel.Delete()

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}
