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
func GetSettingNotification(c *gin.Context) {
	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	accountSettingNotificationModel := &models.AccountSettingNotification{
		AccountUUID: authPayload.AccountUUID,
	}
	result := accountSettingNotificationModel.ReadByAccountUUID()

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
		Data: getSettingNotificationResponsePayload{
			FollowOrOwnArticleReply: accountSettingNotificationModel.FollowOrOwnArticleReply,
			CommentTagged:           accountSettingNotificationModel.CommentTagged,
			ArticleTweet:            accountSettingNotificationModel.ArticleTweet,
			CommentTweet:            accountSettingNotificationModel.CommentTweet,
			InterestRecommendation:  accountSettingNotificationModel.InterestRecommendation,
			Chat:                    accountSettingNotificationModel.Chat,
			Followed:                accountSettingNotificationModel.Followed,
		},
	})
}

// TODO: Doc
func PutSettingNotification(c *gin.Context) {
	var requestPayload *putSettingNotificationRequestPayload

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

	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	m := controllers.GetFilteredNilRequestPayloadMap(requestPayload)
	accountSettingNotificationModel := &models.AccountSettingNotification{
		AccountUUID: authPayload.AccountUUID,
	}
	result := accountSettingNotificationModel.UpdateByAccountUUID(m)
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
