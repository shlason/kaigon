package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
)

func OAuthCallbackForGoogle(c *gin.Context) {

}

const captchaCodeLength int = 6

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
