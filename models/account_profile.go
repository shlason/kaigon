package models

import "gorm.io/gorm"

type AccountProfile struct {
	gorm.Model
	AccountID   uint
	AccountUUID string
	Avatar      string
	Banner      string
	Signature   string
}
