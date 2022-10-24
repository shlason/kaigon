package account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary     註冊帳號
// @Description 使用 Email, Password, Captcha Info 等資訊來註冊
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       email       body     string true "Account Email"
// @Param       password    body     string true "Account Password"
// @Param       captchaUuid body     string true "Captcha Info"
// @Param       captchaCode body     string true "Captcha Info"
// @Success     200         {object} controllers.JSONResponse
// @Failure     400         {object} controllers.JSONResponse
// @Failure     409         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /account/signup [post]
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
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	utils.SendEmail([]string{requestPayload.Email}, "[Kaigon]：恭喜您註冊成功", "signup_success.html", struct{}{})
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

// @Summary     登入帳號
// @Description 使用 Email, Password, Captcha Info 等資訊來登入
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       email       body     string true "Account Email"
// @Param       password    body     string true "Account Password"
// @Param       captchaUuid body     string true "Captcha Info"
// @Param       captchaCode body     string true "Captcha Info"
// @Success     200         {object} controllers.JSONResponse
// @Header      200         {string} Cookie "Refresh Token"
// @Failure     400         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /account/signin [post]
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

	session := &models.Session{
		AccountUUID: accountModel.UUID,
		Email:       accountModel.Email,
	}
	err = session.Read()
	if !(err == nil || err == redis.Nil) {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisGetKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	if err == redis.Nil {
		if session.Create() != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerRedisSetNXKeyGotError,
				Message: err,
				Data:    nil,
			})
			return
		}
	}
	c.SetCookie(
		controllers.RefreshTokenCookieInfo.Name,
		session.Token,
		controllers.RefreshTokenCookieInfo.MaxAge,
		controllers.RefreshTokenCookieInfo.Path,
		controllers.RefreshTokenCookieInfo.Domain,
		controllers.RefreshTokenCookieInfo.Secure,
		controllers.RefreshTokenCookieInfo.HttpOnly,
	)
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

// @Summary     啟動驗證階段
// @Description 發送 Email 認證信
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       accountUUID path     string true "Account UUID"
// @Param       email       body     string true "Account Email"
// @Param       type        body     string true "Verification type"
// @Success     200         {object} controllers.JSONResponse
// @Failure     400         {object} controllers.JSONResponse
// @Failure     401         {object} controllers.JSONResponse
// @Failure     403         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /account/{accountUUID}/info/verification [post]
func CreateVerifySession(c *gin.Context) {
	var requestPayload *createVerifySessionRequestPayload
	errResponse, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	authPayload := c.MustGet("authPayload").(*models.JWTToken)
	// TODO: 以後增加多種驗證方式時，需加上方式的判斷來個別處理
	if authPayload.AccountUUID != c.Param("accountUUID") || authPayload.Email != requestPayload.Email {
		c.JSON(http.StatusForbidden, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPermissionForbidden,
			Message: controllers.ErrMessageRequestPermissionForbidden,
			Data:    nil,
		})
		return
	}

	templatParams := &verificationSessionTemplateParams{}
	err = templatParams.generate(authPayload.AccountUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisSetNXKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	to := []string{
		requestPayload.Email,
	}
	err = utils.SendEmail(to, "[Kaigon]：驗證 Kaigon 所註冊之 Email 的操作指示", "email_verification.html", templatParams)

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

// @Summary     Email 認證信中的認證連結
// @Description 進行 Email 的認證
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       accountUUID path  string true "Account UUID"
// @Param       token       query string true "Session Token"
// @Param       code        query string true "Verify Code"
// @Success     301
// @Failure     400 {object} controllers.JSONResponse
// @Failure     500 {object} controllers.JSONResponse
// @Router      /account/{accountUUID}/info/verification/email [get]
func VerifyWithEmail(c *gin.Context) {
	// TODO: Docs
	requestPayload := &verifyWithEmailRequestPayload{
		AccountUUID: c.Param("accountUUID"),
		Token:       c.Query("token"),
		Code:        c.Query("code"),
	}
	authAccountEmailVerificationModel := &models.AuthAccountEmailVerification{
		AccountUUID: requestPayload.AccountUUID,
		Token:       requestPayload.Token,
		Code:        requestPayload.Code,
	}
	err := authAccountEmailVerificationModel.Read()
	if !(err == nil && authAccountEmailVerificationModel.IsMatch()) {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadTokenCodeFieldsNotValid,
			Message: errMessageRequestPayloadTokenCodeFieldsNotValid,
			Data:    nil,
		})
		return
	}
	accountModel := &models.Account{
		UUID: requestPayload.AccountUUID,
	}
	result := accountModel.UpdateIsEmailVerifiedToTrueByAccountUUID()
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseUpdateGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s://%s", configs.Server.Protocol, configs.Server.Host))
}

// @Summary     啟動重設密碼階段 (忘記密碼時)
// @Description 發送可以重設密碼的連結到指定 email 中
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       email body     string true "Account Email"
// @Param       path  body     string true "Front End Reset Password URL's Path"
// @Success     200   {object} controllers.JSONResponse
// @Failure     400   {object} controllers.JSONResponse
// @Failure     500   {object} controllers.JSONResponse
// @Router      /account/info/password/reset [post]
func CreateResetPasswordSession(c *gin.Context) {
	var requestPayload *createResetPasswordSessionRequestPayload
	errResponse, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
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
	err = templatParams.generate(accountModel.UUID, requestPayload.Email, requestPayload.Path)

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

// @Summary     強制重設密碼 (忘記密碼時)
// @Description 強制重新設定密碼藉由 Token, Code 相關驗證資訊
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       email       body     string true "Account Email"
// @Param       password    body     string true "Account New Password"
// @Param       token       body     string true "Token"
// @Param       code        body     string true "Verify Code"
// @Param       captchaUuid body     string true "Captcha Info"
// @Param       captchaCode body     string true "Captcha Info"
// @Success     200         {object} controllers.JSONResponse
// @Failure     400         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /account/info/password/reset [patch]
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
