package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Email           string `gorm:"unique; not null;"`
	Password        string
	IsEmailVerified bool
}
