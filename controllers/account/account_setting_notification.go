package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"gorm.io/gorm"
)

// @Summary     取得 Account 通知設定
// @Description 取得 Account 相關通知設定
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} controllers.JSONResponse{Data=getSettingNotificationResponsePayload}
// @Failure     400 {object} controllers.JSONResponse
// @Failure     500 {object} controllers.JSONResponse
// @Router      /account/:accountUUID/setting/notification [get]
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

// @Summary     更新 Account 通知設定
// @Description 更新 Account 相關通知設定
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       followOrOwnArticleReply body     uint true "我追蹤的/我的文章有新留言 (0: 關閉, 1: 所有留言, 2: 標註我的留言)"
// @Param       commentTagged           body     bool true "我的留言被標注"
// @Param       articleTweet            body     bool true "我的文章獲得心情"
// @Param       commentTweet            body     bool true "我的留言獲得愛心"
// @Param       interestRecommendation  body     bool true "我可能感興趣的內容"
// @Param       chat                    body     bool true "聊天通知"
// @Param       followed                body     bool true "被人追蹤時"
// @Success     200                     {object} controllers.JSONResponse
// @Failure     400                     {object} controllers.JSONResponse
// @Failure     500                     {object} controllers.JSONResponse
// @Router      /account/:accountUUID/setting/notification [put]
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
