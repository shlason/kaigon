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

	accountProfileModel := &models.AccountProfile{
		AccountID:   accountModel.ID,
		AccountUUID: accountModel.UUID,
	}
	result = accountProfileModel.Create()
	if result.Error != nil {
		return JSONResponse{
			Code:    ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		}, true
	}

	accountSettingModel := &models.AccountSetting{
		AccountID:   accountModel.ID,
		AccountUUID: accountModel.UUID,
	}
	result = accountSettingModel.Create()
	if result.Error != nil {
		return JSONResponse{
			Code:    ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		}, true
	}

	accountSettingNotificationModel := &models.AccountSettingNotification{
		AccountID:   accountModel.ID,
		AccountUUID: accountModel.UUID,
	}
	result = accountSettingNotificationModel.Create()
	if result.Error != nil {
		return JSONResponse{
			Code:    ErrCodeServerDatabaseCreateGotError,
			Message: result.Error,
			Data:    nil,
		}, true
	}

	return JSONResponse{}, false
}
