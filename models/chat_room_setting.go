package models

import (
	"gorm.io/gorm"
)

type ChatRoomSetting struct {
	gorm.Model
	ChatRoomID uint `gorm:"unique; not null;"`
	Emoji      string
	Name       string
	Avatar     string
}

func (crs *ChatRoomSetting) Create() *gorm.DB {
	return db.Create(&crs)
}
