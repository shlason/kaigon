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
	"github.com/shlason/kaigon/models/constants"
	"github.com/shlason/kaigon/utils"
	"gorm.io/gorm"
)

const captchaCodeLength int = 8

// @Summary     取得 Google OAuth 相關 URL
// @Description 帶上 type (login, bind) 來表示要取得的是登入用，還是綁定第三方登入用，以及 redirectPath 表示成功後所導轉的終點
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       type         query    string true "login or bind"
// @Param       redirectPath query    string true "導轉終點"
// @Success     200          {object} controllers.JSONResponse{data=getOAuthUrlResponsePayload}
// @Failure     400          {object} controllers.JSONResponse
// @Router      /auth/o/google/url [get]
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
				"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&state=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.email",
				configs.OAuth.Google.ClientID,
				requestParams.RedirectPath,
				fmt.Sprintf(
					"%s://%s/api/auth/o/google/%s",
					configs.Server.Protocol,
					configs.Server.Host,
					requestParams.Type,
				),
			),
		},
	})
}

func getGoogleOAuthInfoByGrantCodeAndURLType(c *gin.Context, grantCode string, URLType string) (googleOAuthUserInfoResponsePayload, error) {
	requestPayload, err := json.Marshal(map[string]string{
		"client_id":     configs.OAuth.Google.ClientID,
		"client_secret": configs.OAuth.Google.ClientSecret,
		"code":          grantCode,
		"redirect_uri":  fmt.Sprintf("%s://%s/api/auth/o/google/%s", configs.Server.Protocol, configs.Server.Host, URLType),
		"grant_type":    "authorization_code",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	responseBody := bytes.NewBuffer(requestPayload)

	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", responseBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    ErrCodeRequestOAuthAccessTokenGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	accessTokenResp := googleOAuthAccessTokenResponsePayload{}
	err = json.Unmarshal(body, &accessTokenResp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", accessTokenResp.TokenType, accessTokenResp.AccessToken))
	fresp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    ErrCodeRequestOAuthAccessTokenGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	defer fresp.Body.Close()

	body, err = ioutil.ReadAll(fresp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	userInfoPayload := googleOAuthUserInfoResponsePayload{}
	err = json.Unmarshal(body, &userInfoPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return googleOAuthUserInfoResponsePayload{}, err
	}

	return userInfoPayload, nil
}

// TODO: 登入、註冊有很多重複的 CODE 需要整理
// @Summary     Google OAuth 登入且授權成功後所導轉的 URI (登入用)
// @Description 拿到 Google 給的 grant code 後再去和 Google 拿 access token，再用 access token 去拿該 user 的相關資訊 (scope)，若失敗會在 QS 加上 oauth_login_failed=1
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       code  query string true "Google grant code (from google)"
// @Param       state query string true "Front end redirect path (from google OAuth state QS)"
// @Success     302
// @Failure     302
// @Failure     400 {object} controllers.JSONResponse
// @Failure     500 {object} controllers.JSONResponse
// @Router      /auth/o/google/login [get]
func GoogleOAuthRedirectURIForLogin(c *gin.Context) {
	var requestParams googleOAuthRedirectURIForLoginQueryParmas

	err := c.ShouldBindQuery(&requestParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestQueryParamsNotValid,
			Message: controllers.ErrMessageRequestQueryParamsNotValid,
			Data:    nil,
		})
		return
	}

	userInfoPayload, err := getGoogleOAuthInfoByGrantCodeAndURLType(c, requestParams.Code, "login")

	if err != nil {
		return
	}

	accountOAuthInfoModel := &models.AccountOauthInfo{
		Email:    userInfoPayload.Email,
		Provider: OAuthProviderName["google"],
	}

	// 是否有綁定過第三方登入方式
	result := accountOAuthInfoModel.ReadByEmailAndProvider()

	// 綁定過第三方登入方式的話直接去 account model 查詢來直接登入
	if result.Error == nil {
		am := &models.Account{
			UUID: accountOAuthInfoModel.AccountUUID,
		}
		r := am.ReadByUUID()
		// 不管是查詢不到紀錄或不明錯誤都是有問題，一律視作 500
		if r.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseQueryGotError,
				Message: err,
				Data:    nil,
			})
			return
		}
		// 執行登入
		session := &models.Session{
			AccountID:   am.ID,
			AccountUUID: am.UUID,
			Email:       am.Email,
		}
		err = session.ReadByAccount()

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
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     constants.RefreshTokenCookieInfo.Name,
			Value:    session.Token,
			Path:     constants.RefreshTokenCookieInfo.Path,
			Domain:   constants.RefreshTokenCookieInfo.Domain,
			MaxAge:   constants.RefreshTokenCookieInfo.MaxAge,
			Secure:   constants.RefreshTokenCookieInfo.Secure,
			HttpOnly: constants.RefreshTokenCookieInfo.HttpOnly,
			SameSite: http.SameSite(constants.RefreshTokenCookieInfo.SameSite),
		})
		c.Redirect(http.StatusFound, fmt.Sprintf(
			"%s://%s%s?success=1",
			configs.Server.Protocol,
			configs.Server.Host,
			requestParams.State,
		))
		return
	}

	// 若該第三方登入方式還未綁定過，則去檢查該第三方的 email 是否已存在於 account model 中
	// 若不存在就使用該 email 註冊一個帳號並且與該第三方登入方式關聯起來
	// 若存在則跳錯誤告知使用者說該 email 已被註冊過，可能可以單純使用 email 方式登入，後續進後台在再綁定該第三方登入方式
	accountModel := &models.Account{
		UUID:  uuid.NewString(),
		Email: userInfoPayload.Email,
	}
	result = accountModel.ReadByEmail()
	// TODO: Error handle email 存在時的情境
	// TODO: 和前端討論 失敗後導轉回去時要帶的 Query
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Redirect(http.StatusFound, fmt.Sprintf(
			"%s://%s%s?oauth_login_failed=1",
			configs.Server.Protocol,
			configs.Server.Host,
			requestParams.State,
		))
		return
	}
	// Email 不存在則自動註冊一個新帳號並關聯第三方登入
	errResp, hasErr := controllers.InitAccountDataWhenSignUp(accountModel)
	if hasErr {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	accountOAuthInfoModel.AccountID = accountModel.ID
	accountOAuthInfoModel.AccountUUID = accountModel.UUID

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

	// 執行登入
	session := &models.Session{
		AccountID:   accountModel.ID,
		AccountUUID: accountModel.UUID,
		Email:       accountModel.Email,
	}

	err = session.ReadByAccount()
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshTokenCookieInfo.Name,
		Value:    session.Token,
		Path:     constants.RefreshTokenCookieInfo.Path,
		Domain:   constants.RefreshTokenCookieInfo.Domain,
		MaxAge:   constants.RefreshTokenCookieInfo.MaxAge,
		Secure:   constants.RefreshTokenCookieInfo.Secure,
		HttpOnly: constants.RefreshTokenCookieInfo.HttpOnly,
		SameSite: http.SameSite(constants.RefreshTokenCookieInfo.SameSite),
	})

	c.Redirect(http.StatusFound, fmt.Sprintf(
		"%s://%s%s?success=1",
		configs.Server.Protocol,
		configs.Server.Host,
		requestParams.State,
	))
}

// @Summary     Google OAuth 登入且授權成功後所導轉的 URI (綁定用)
// @Description 拿到 Google 給的 grant code 後導轉去指定前端頁面路徑 (導轉路徑為 "[GET]取得 Google OAuth 相關 URL" 當初打 API 時給什麼路徑就轉去那)
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       code  query string true "Google grant code (from google)"
// @Param       state query string true "Front end redirect path (from google OAuth state QS)"
// @Success     302
// @Failure     400 {object} controllers.JSONResponse
// @Router      /auth/o/google/bind [get]
func GoogleOAuthRedirectURIForBind(c *gin.Context) {
	var requestParams *googleOAuthRedirectURIForBindQueryParams

	err := c.ShouldBindQuery(&requestParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestQueryParamsNotValid,
			Message: controllers.ErrMessageRequestQueryParamsNotValid,
			Data:    nil,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf(
		"%s://%s%s?grantCode=%s",
		configs.Server.Protocol,
		configs.Server.Host,
		requestParams.State,
		requestParams.Code,
	))
}

// @Summary     更新 Account 第三方登入綁定資訊
// @Description 更新 Account 第三方登入綁定資訊
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       grantCode body     string true "Google grant code (from google)"
// @Success     200       {object} controllers.JSONResponse
// @Failure     400       {object} controllers.JSONResponse
// @Failure     409       {object} controllers.JSONResponse
// @Failure     500       {object} controllers.JSONResponse
// @Router      /auth/o/google/bind [patch]
func GoogleOAuthBind(c *gin.Context) {
	var requestPayload *googleOAuthBindRequestPayload

	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	userInfoPayload, err := getGoogleOAuthInfoByGrantCodeAndURLType(c, requestPayload.GrantCode, "bind")

	if err != nil {
		return
	}

	accountOAuthInfoModel := &models.AccountOauthInfo{
		Email:    userInfoPayload.Email,
		Provider: OAuthProviderName["google"],
	}

	// 是否有綁定過第三方登入方式
	result := accountOAuthInfoModel.ReadByEmailAndProvider()

	// DB query 沒遇到問題 (record not found) 代表已經有該 OAuth Info
	if result.Error == nil {
		c.JSON(http.StatusConflict, controllers.JSONResponse{
			Code:    ErrCodeRequestRecordAlreadyExistWhenOAuthBinding,
			Message: ErrMessageRequestRecordAlreadyExistWhenOAuthBinding,
			Data:    nil,
		})
		return
	}
	// 執行到這裡代表 DB query 遇到問題，若問題不是預期的 (record not found) 則當作不明錯誤 直接噴 500
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	accountOAuthInfoModel.AccountID = authPayload.AccountID
	accountOAuthInfoModel.AccountUUID = authPayload.AccountUUID

	result = accountOAuthInfoModel.Create()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
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

// @Summary     取得 authToken
// @Description 使用 Cookie 中的 REFRESH_TOKEN field 來獲取 authToken
// @Tags        auth
// @Accept      json
// @Produce     json
// @Success     200 {object} controllers.JSONResponse{data=getAuthTokenByRefreshTokenResponsePayload}
// @Failure     400 {object} controllers.JSONResponse
// @Failure     401 {object} controllers.JSONResponse
// @Failure     500 {object} controllers.JSONResponse
// @Router      /auth/session/token/refresh [get]
func GetAuthTokenByRefreshToken(c *gin.Context) {
	token, err := c.Cookie(constants.RefreshTokenCookieInfo.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Message: ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Data:    nil,
		})
		return
	}

	session := models.Session{
		Token: token,
	}
	err = session.ReadByToken()

	if err == redis.Nil {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     constants.RefreshTokenCookieInfo.Name,
			Value:    "",
			Path:     constants.RefreshTokenCookieInfo.Path,
			Domain:   constants.RefreshTokenCookieInfo.Domain,
			MaxAge:   -1,
			Secure:   constants.RefreshTokenCookieInfo.Secure,
			HttpOnly: constants.RefreshTokenCookieInfo.HttpOnly,
			SameSite: http.SameSite(constants.RefreshTokenCookieInfo.SameSite),
		})
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Message: ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized,
			Data:    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerRedisGetKeyGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	jwtModel := &models.JWTToken{
		AccountID:   session.AccountID,
		AccountUUID: session.AccountUUID,
		Email:       session.Email,
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
