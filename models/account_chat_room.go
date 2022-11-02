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

func (acr *AccountChatRoom) Create() *gorm.DB {
	return db.Create(&acr)
}

func (acr AccountChatRoom) ReadAllByChatRoomID(chatRoomID uint, list *[]AccountChatRoom) *gorm.DB {
	return db.Where("chat_room_id = ?", chatRoomID).Find(&list)
}
