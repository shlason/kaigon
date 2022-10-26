package controllers

import "github.com/shlason/kaigon/models"

func InitAccountDataWhenSignUp(accountModel *models.Account) (errResponse JSONResponse, hasErr bool) {
	result := accountModel.Create()

	if result.Error != nil {
		return JSONResponse{
			Code:    ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		}, true
	}

	return JSONResponse{}, false
}
