package models

import "gorm.io/gorm"

type AccountOauthInfo struct {
	gorm.Model
	AccoundID   uint
	AccountUUID string
	Provider    string
	Email       string
}

func (accountOAuthInfo *AccountOauthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}
