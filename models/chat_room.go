package models

import (
	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	Type             string
	MaximumMemberNum int
}
