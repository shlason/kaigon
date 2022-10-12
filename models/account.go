package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Uuid            string `gorm:"unique; not null;"`
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

func (account *Account) UpdatePasswordByEmail(password string) *gorm.DB {
	return db.Model(&account).Where("email = ?", account.Email).Update("password", password)
}
