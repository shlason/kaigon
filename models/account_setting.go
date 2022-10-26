package models

import "gorm.io/gorm"

type AccountSetting struct {
	gorm.Model
	AccountID   uint   `gorm:"unique; not null;"`
	AccountUUID string `gorm:"unique; not null;"`
	Name        string
	Locale      string
}

func (accountSetting *AccountSetting) Create() *gorm.DB {
	return db.Create(&accountSetting)
}
