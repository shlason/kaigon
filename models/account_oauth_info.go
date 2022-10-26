package models

import "gorm.io/gorm"

type AccountOauthInfo struct {
	gorm.Model
	AccountID   uint   `gorm:"not null"`
	AccountUUID string `gorm:"not null"`
	Provider    string `gorm:"not null"`
	Email       string `gorm:"not null"`
}

func (accountOAuthInfo *AccountOauthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}

func (accountOAuthInfo *AccountOauthInfo) ReadByEmailAndProvider() *gorm.DB {
	return db.First(&accountOAuthInfo, "email = ? AND provider = ?", accountOAuthInfo.Email, accountOAuthInfo.Provider)
}
