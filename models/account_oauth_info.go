package models

import "gorm.io/gorm"

type AccountOAuthInfo struct {
	gorm.Model
	AccoundID   uint
	AccountUUID string
	Provider    string
	Email       string
}

func (accountOAuthInfo *AccountOAuthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}
