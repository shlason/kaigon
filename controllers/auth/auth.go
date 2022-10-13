package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/utils"
)

func OAuthCallbackForGoogle(c *gin.Context) {

}

func GetCaptchaInfo(c *gin.Context) {
	const captchaCodeLength int = 6

	authCaptchaModel := &models.AuthCaptcha{
		UUID: uuid.NewString(),
		Code: utils.RandStringBytes(captchaCodeLength),
	}
	err := authCaptchaModel.Create()
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
		Data: getCaptchaInfoResponsePayload{
			UUID: authCaptchaModel.UUID,
		},
	})
}

func GetCaptchaImage(c *gin.Context) {

	fmt.Println("Suss")
	//c.Data(http.StatusOK, "image/png", buffer.Bytes())
}

func UpdateCaptchaInfo(c *gin.Context) {

}
