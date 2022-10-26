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
func GetProfile(c *gin.Context) {
	accountProfileModel := &models.AccountProfile{
		AccountUUID: c.Param("accountUUID"),
	}
	result := accountProfileModel.ReadByAccountUUID()

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

	var socialMedias []models.AccountProfileSocialMedia
	var socialMediasResponsePayload []socialMediaResponsePayload

	result = models.AccountProfileSocialMedia{}.ReadAllByAccountUUID(c.Param("accountUUIID"), &socialMedias)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: result.Error,
			Data:    nil,
		})
		return
	}

	for _, socialMedia := range socialMedias {
		socialMediasResponsePayload = append(socialMediasResponsePayload, socialMediaResponsePayload{
			Provider: socialMedia.Provider,
			UserName: socialMedia.UserName,
		})
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: getProfileResponsePayload{
			Avatar:       accountProfileModel.Avatar,
			Banner:       accountProfileModel.Banner,
			Signature:    accountProfileModel.Signature,
			SocialMedias: socialMediasResponsePayload,
		},
	})
}

// TODO: Doc
func PatchProfile(c *gin.Context) {
	var requestPayload *patchProfileRequestPayload
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

	m := controllers.GetFilteredNilRequestPayloadMap(requestPayload)

	authPayload := c.MustGet("authPayload").(*models.JWTToken)
	accountProfileModel := &models.AccountProfile{
		AccountUUID: authPayload.AccountUUID,
	}
	result := accountProfileModel.UpdateByAccountUUID(m)
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
		Data:    nil,
	})
}
