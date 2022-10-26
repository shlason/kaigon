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

func (accountProfile *AccountProfile) ReadByAccountUUID() *gorm.DB {
	return db.First(&accountProfile, "account_uuid = ?", accountProfile.AccountUUID)
}

func (accountProfile *AccountProfile) UpdateByAccountUUID(m map[string]interface{}) *gorm.DB {
	return db.Model(&accountProfile).Where("account_uuid = ?", accountProfile.AccountUUID).Updates(m)
}
