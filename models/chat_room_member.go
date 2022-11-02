package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoomMember struct {
	ID                  uint   `gorm:"primarykey"`
	ChatRoomID          uint   `gorm:"not null;"`
	AccountUUID         string `gorm:"not null;"`
	Theme               string
	EnabledNotification bool `gorm:"default:true; not null;"`
	LastSeenAt          time.Time
}

func (acr *ChatRoomMember) Create() *gorm.DB {
	return db.Create(&acr)
}

func (acr ChatRoomMember) ReadAllByChatRoomID(chatRoomID uint, list *[]ChatRoomMember) *gorm.DB {
	return db.Where("chat_room_id = ?", chatRoomID).Find(&list)
}

func (acr ChatRoomMember) ReadAllByAccountUUID(accountUUID string, list *[]ChatRoomMember) *gorm.DB {
	return db.Where("account_uuid = ?", accountUUID).Find(&list)
}
