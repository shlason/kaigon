package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UUID            string `gorm:"unique; not null;"`
	Email           string `gorm:"unique; not null;"`
	Password        string
	IsEmailVerified bool
}

func (account *Account) Create() *gorm.DB {
	return db.Create(&account)
}

func (account *Account) ReadByEmail() *gorm.DB {
	return db.First(&account, "email = ?", account.Email)
}

func (account *Account) ReadByUUID() *gorm.DB {
	return db.First(&account, "uuid = ?", account.UUID)
}

func (account *Account) UpdateByEmailAndUUID(m map[string]interface{}) *gorm.DB {
	return db.Model(&account).Where("uuid = ? AND email = ?", account.UUID, account.Email).Updates(m)
}

func (account *Account) UpdatePasswordByEmail(password string) *gorm.DB {
	return db.Model(&account).Where("email = ?", account.Email).Update("password", password)
}

func (account *Account) UpdateIsEmailVerifiedToTrueByAccountUUID() *gorm.DB {
	return db.Model(&account).Where("uuid = ?", account.UUID).Update("is_email_verifieed", true)
}
