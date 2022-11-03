package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
)

func RegisteDevelopUtilsRoutes(r *gin.RouterGroup) {
	// @Summary     刪除所有開發者的所有帳號資料
	// @Description 刪除所有開發者的所有帳號資料 (直接在網頁爆打這支就可以刪掉資料了)
	// @Tags        useful_utils_when_developing
	// @Accept      json
	// @Produce     json
	// @Success     200 {object} controllers.JSONResponse
	// @Failure     500 {object} controllers.JSONResponse
	// @Router      /develop/utils/account/delete [get]
	r.GET("/develop/utils/account/delete", func(c *gin.Context) {
		var hardcodeEmails = []interface{}{
			"nocvi111@gmail.com",
			"hn15637648@yahoo.com.tw",
			"hn15637648@gmail.com",
			"a7636439@gmail.com",
			"yanalu0818@gmail.com",
		}

		var accountModels []models.Account

		result := models.Account{}.ReadByEmails(hardcodeEmails, &accountModels)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseQueryGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		var accountIDs []interface{}

		for _, account := range accountModels {
			accountIDs = append(accountIDs, account.ID)
		}

		result = models.Account{}.DeleteByIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		result = models.AccountProfile{}.DeleteByAccountIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		result = models.AccountSetting{}.DeleteByAccountIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		result = models.AccountSettingNotification{}.DeleteByAccountIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		result = models.AccountOauthInfo{}.DeleteByAccountIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
				Message: result.Error,
				Data:    nil,
			})
			return
		}

		result = models.AccountProfileSocialMedia{}.DeleteByAccountIDs(accountIDs)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
				Code:    controllers.ErrCodeServerDatabaseDeleteGotError,
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
	})
}
