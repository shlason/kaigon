package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

func GetSetting(c *gin.Context) {
	accountSettingModel := models.AccountSetting{
		AccountUUID: c.Param("accountUUID"),
	}

	result := accountSettingModel.ReadByAccountUUID()

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, controllers.JSONResponse{
				Code:    errCodeRequestParamAccountUUIDNotFound,
				Message: errMessageRequestParamAccountUUIDNotFound,
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

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getSettingResponsePayload{
			GormModelResponse: controllers.GormModelResponse{
				CreatedAt: accountSettingModel.CreatedAt,
				UpdatedAt: accountSettingModel.UpdatedAt,
				DeletedAt: accountSettingModel.DeletedAt,
			},
			Name:   accountSettingModel.Name,
			Locale: accountSettingModel.Locale,
		},
	})
}

func PatchSetting(c *gin.Context) {

}
