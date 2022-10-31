package models

import (
	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	Type             string `gorm:"not null;"`
	MaximumMemberNum int    `gorm:"not null;"`
}

func (cr *ChatRoom) Create() *gorm.DB {
	return db.Create(&cr)
}
