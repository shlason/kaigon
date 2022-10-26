package models

import "gorm.io/gorm"

type AccountProfile struct {
	gorm.Model
	AccountID   uint   `gorm:"unique; not null;"`
	AccountUUID string `gorm:"unique; not null;"`
	Avatar      string
	Banner      string
	Signature   string
}

func (accountProfile *AccountProfile) Create() *gorm.DB {
	return db.Create(&accountProfile)
}
