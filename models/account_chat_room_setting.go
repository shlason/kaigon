package models

import (
	"gorm.io/gorm"
)

type AccountChatRoomSetting struct {
	gorm.Model
	AccountChatRoomID   uint `gorm:"unique; not null;"`
	Theme               string
	EnabledNotification bool `gorm:"default:true; not null;"`
}
