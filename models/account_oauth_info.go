package models

import "gorm.io/gorm"

type AccountOauthInfo struct {
	gorm.Model
	AccoundID   uint   `gorm:"unique; not null;"`
	AccountUUID string `gorm:"unique; not null;"`
	Provider    string
	Email       string
}

func (accountOAuthInfo *AccountOauthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}

func (accountOAuthInfo *AccountOauthInfo) ReadByEmailAndProvider() *gorm.DB {
	return db.First(&accountOAuthInfo, "email = ? AND provider = ?", accountOAuthInfo.Email, accountOAuthInfo.Provider)
}
