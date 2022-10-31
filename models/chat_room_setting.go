package models

import (
	"gorm.io/gorm"
)

type ChatRoomSetting struct {
	gorm.Model
	ChatRoomID uint
	Emoji      string
	Name       string
	Avatar     string
}
