package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
)

func RegisteDevelopUtilsRoutes(r *gin.RouterGroup) {
	r.GET("/develop/utils/account/delete", func(c *gin.Context) {
		var hardcodeEmails = []interface{}{
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

	})
}
