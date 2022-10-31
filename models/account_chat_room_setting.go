package models

import (
	"gorm.io/gorm"
)

type AccountChatRoomSetting struct {
	gorm.Model
	AccountChatRoomID   uint
	Theme               string
	EnabledNotification bool
}
