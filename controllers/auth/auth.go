package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
	"gorm.io/gorm"
)

const captchaCodeLength int = 6

// TODO: Doc
func GetGoogleOAuthURL(c *gin.Context) {
	var requestParams *getOAuthUrlQueryParmas

	err := c.ShouldBindQuery(&requestParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestQueryParamsNotValid,
			Message: controllers.ErrMessageRequestQueryParamsNotValid,
			Data:    nil,
		})
		return
	}

	errResp, isNotValid := requestParams.check()

	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getOAuthUrlResponsePayload{
			URL: fmt.Sprintf(
				"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.email",
				configs.OAuth.Google.ClientID,
				fmt.Sprintf("%s://%s/api/auth/o/google/%s", configs.Server.Protocol, configs.Server.Host, requestParams.Type),
			),
		},
	})
}

// TODO: Doc
func GoogleOAuthRedirectURIForLogin(c *gin.Context) {
	requestPayload, err := json.Marshal(map[string]string{
		"client_id":     configs.OAuth.Google.ClientID,
		"client_secret": configs.OAuth.Google.ClientSecret,
		"code":          c.Query("code"),
		"redirect_uri":  fmt.Sprintf("%s://%s/api/auth/o/google/login", configs.Server.Protocol, configs.Server.Host),
		"grant_type":    "authorization_code",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	responseBody := bytes.NewBuffer(requestPayload)

	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", responseBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    ErrCodeRequestOAuthAccessTokenGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	accessTokenResp := googleOAuthAccessTokenResponsePayload{}
	err = json.Unmarshal(body, &accessTokenResp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", accessTokenResp.TokenType, accessTokenResp.AccessToken))
	fresp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    ErrCodeRequestOAuthAccessTokenGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	defer fresp.Body.Close()

	body, err = ioutil.ReadAll(fresp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	userInfoPayload := googleOAuthUserInfoResponsePayload{}
	err = json.Unmarshal(body, &userInfoPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	accountModel := &models.Account{
		UUID:  uuid.NewString(),
		Email: userInfoPayload.Email,
	}
	result := accountModel.ReadByEmail()
	// TODO: Error handle
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, controllers.JSONResponse{
			Code:    "0",
			Message: "m",
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

	accountOAuthInfoModel := &models.AccountOAuthInfo{
		AccoundID:   accountModel.ID,
		AccountUUID: accountModel.UUID,
		Email:       accountModel.Email,
		Provider:    OAuthProviderName["google"],
	}

	result = accountOAuthInfoModel.Create()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	utils.SendEmail([]string{accountModel.Email}, "[Kaigon]：恭喜您註冊成功", "signup_success.html", struct{}{})

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}

// @Summary     取得 authToken
// @Description 使用 Cookie 中的 REFRESH_TOKEN field 來獲取 authToken
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       accountUuid query    string true "Account UUID"
// @Param       email       query    string true "Account Email"
// @Success     200         {object} controllers.JSONResponse{data=getAuthTokenByRefreshTokenResponsePayload}
// @Failure     400         {object} controllers.JSONResponse
// @Failure     401         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /auth/session/token/refresh [get]
func GetAuthTokenByRefreshToken(c *gin.Context) {
	token, err := c.Cookie(controllers.RefreshTokenCookieInfo.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Message: ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Data:    nil,
		})
		return
	}
	requestParams := getAuthTokenByRefreshTokenRequestParamsPayload{}
	err = c.BindQuery(&requestParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    ErrCodeRequestQueryParamsAccountUUIDFieldNotValid,
			Message: ErrMessageRequestQueryParamsAccountUUIDFieldNotValid,
			Data:    nil,
		})
		return
	}
	errResp, isNotValid := requestParams.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	session := models.Session{
		AccountUUID: requestParams.AccountUUID,
		Email:       requestParams.Email,
	}
	err = session.Read()
	if err != nil {
		if err == redis.Nil {
			c.SetCookie(
				controllers.RefreshTokenCookieInfo.Name,
				"",
				-1,
				controllers.RefreshTokenCookieInfo.Path,
				controllers.RefreshTokenCookieInfo.Domain,
				controllers.RefreshTokenCookieInfo.Secure,
				controllers.RefreshTokenCookieInfo.HttpOnly,
			)
			c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
				Code:    ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized,
				Message: ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisGetKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	if session.Token != token {
		c.SetCookie(
			controllers.RefreshTokenCookieInfo.Name,
			"",
			-1,
			controllers.RefreshTokenCookieInfo.Path,
			controllers.RefreshTokenCookieInfo.Domain,
			controllers.RefreshTokenCookieInfo.Secure,
			controllers.RefreshTokenCookieInfo.HttpOnly,
		)
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Message: ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Data:    nil,
		})
		return
	}

	jwtModel := &models.JWTToken{
		AccountUUID: requestParams.AccountUUID,
		Email:       requestParams.Email,
	}

	authToken, err := jwtModel.Generate()

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGenerateJWTTokenGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getAuthTokenByRefreshTokenResponsePayload{
			AuthToken: authToken,
		},
	})
}

// @Summary     取得圖形驗證相關資訊
// @Description 取得圖形驗證動態產生的相對應 UUID 及驗證圖片
// @Tags        auth
// @Accept      json
// @Produce     json
// @Success     200 {object} controllers.JSONResponse{data=getCaptchaInfoResponsePayload}
// @Failure     500 {object} controllers.JSONResponse
// @Router      /auth/captcha [get]
func GetCaptchaInfo(c *gin.Context) {
	authCaptchaModel := &models.AuthCaptcha{
		UUID: uuid.NewString(),
		Code: utils.RandStringBytes(captchaCodeLength),
	}
	err := authCaptchaModel.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisSetNXKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getCaptchaInfoResponsePayload{
			UUID: authCaptchaModel.UUID,
		},
	})
}

// @Summary     取得圖形驗證圖片
// @Description 取得與 UUID 相對應的驗證圖片
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       captchaUUID path     string       true "Captcha Info"
// @Success     200         {string} content-type "image/png"
// @Failure     500         {object} controllers.JSONResponse
// @Router      /auth/captcha/{captchaUUID}/image [get]
func GetCaptchaImage(c *gin.Context) {
	requestParams := &getCaptchaImageRequestParamsPayload{}
	requestParams.BindParams(c)
	authCaptchaModel := &models.AuthCaptcha{
		UUID: requestParams.CaptchaUUID,
	}
	err := authCaptchaModel.ReadByUUID()

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisGetKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	buffer, err := utils.CreateCaptchaImage(authCaptchaModel.Code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
	}

	c.Data(http.StatusOK, "image/png", buffer.Bytes())
}

// @Summary     刷新圖形驗證
// @Description 刷新與 UUID 相對應的驗證圖片資訊及效期
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       captchaUUID path     string true "Captcha Info"
// @Success     200         {object} controllers.JSONResponse
// @Failure     500         {object} controllers.JSONResponse
// @Router      /auth/captcha/{captchaUUID}/refresh [get]
func UpdateCaptchaInfo(c *gin.Context) {
	requestParams := &updateCaptchaImageRequestParamsPayload{}
	requestParams.BindParams(c)
	authCaptchaModel := &models.AuthCaptcha{
		UUID: requestParams.CaptchaUUID,
		Code: utils.RandStringBytes(captchaCodeLength),
	}
	err := authCaptchaModel.UpdateByUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisSetKeyGotError,
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
