package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountChatRoom struct {
	gorm.Model
	ChatRoomID  uint   `gorm:"not null;"`
	AccountUUID string `gorm:"not null;"`
	LastSeenAt  time.Time
}
