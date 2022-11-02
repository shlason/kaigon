package models

import (
	"gorm.io/gorm"
)

type ChatRoomMemberSetting struct {
	gorm.Model
	ChatRoomMemberID    uint `gorm:"unique; not null;"`
	Theme               string
	EnabledNotification bool `gorm:"default:true; not null;"`
}

func (acrs *ChatRoomMemberSetting) Create() *gorm.DB {
	return db.Create(&acrs)
}
