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

}
