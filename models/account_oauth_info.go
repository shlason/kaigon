package models

import "gorm.io/gorm"

type AccountOAuthInfo struct {
	gorm.Model
	AccoundID int
	Provider  string
	Email     string
}

func (accountOAuthInfo *AccountOAuthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}
