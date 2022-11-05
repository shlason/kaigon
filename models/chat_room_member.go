package models

import (
	"strings"
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

func (acr ChatRoomMember) ReadAllByChatRoomIDs(ids []interface{}, list *[]ChatRoomMember) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "chat_room_id = ?")
	}

	return db.Where(strings.Join(fields, " OR "), ids...).Find(&list)
}

func (acr *ChatRoomMember) UpdateByChatRoomIDAndAccountUUID(m map[string]interface{}) *gorm.DB {
	return db.Model(&acr).Where("chat_room_id = ? AND account_uuid = ?", acr.ChatRoomID, acr.AccountUUID).Updates(m)
}
