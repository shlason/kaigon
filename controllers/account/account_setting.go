package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

// TODO: Doc
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

	var accountOAuthInfos []models.AccountOauthInfo
	accountOAuthInfosResponsePayload := []oauthInfoResponsePayload{}

	result = models.AccountOauthInfo{}.ReadAllByAccountUUID(c.Param("accountUUID"), &accountOAuthInfos)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	for _, accountOAuthInfo := range accountOAuthInfos {
		accountOAuthInfosResponsePayload = append(accountOAuthInfosResponsePayload, oauthInfoResponsePayload{
			GormModelResponse: controllers.GormModelResponse{
				CreatedAt: accountOAuthInfo.CreatedAt,
				UpdatedAt: accountOAuthInfo.UpdatedAt,
				DeletedAt: accountOAuthInfo.DeletedAt,
			},
			Email:    accountOAuthInfo.Email,
			Provider: accountOAuthInfo.Provider,
		})
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
			Name:       accountSettingModel.Name,
			Locale:     accountSettingModel.Locale,
			OAuthInfos: accountOAuthInfosResponsePayload,
		},
	})
}

// TODO: Doc
func PatchSetting(c *gin.Context) {
	var requestPayload *patchSettingRequestPayload

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

	m := controllers.GetFilteredNilRequestPayloadMap(&requestPayload)
	authPayload := c.MustGet("authPayload").(*models.JWTToken)
	accountSettingModel := models.AccountSetting{
		AccountUUID: authPayload.AccountUUID,
	}
	result := accountSettingModel.UpdateByAccountUUID(m)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseUpdateGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: patchSettingResponsePayload{
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
