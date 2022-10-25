package models

import "gorm.io/gorm"

type AccountSetting struct {
	gorm.Model
	AccountID   uint
	AccountUUID string
	Name        string
	Locale      string
}
