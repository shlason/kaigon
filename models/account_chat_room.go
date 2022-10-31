package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountChatRoom struct {
	gorm.Model
	AccountUUID string
	ChatRoomID  uint
	LastSeenAt  time.Time
}
