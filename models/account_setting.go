package models

import "gorm.io/gorm"

type AccountSetting struct {
	gorm.Model
	AccountID   uint
	AccountUUID string
	Name        string
	Locale      string
}

func (accountSetting *AccountSetting) Create() *gorm.DB {
	return db.Create(&accountSetting)
}
