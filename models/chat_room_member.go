package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoomMember struct {
	gorm.Model
	ChatRoomID  uint   `gorm:"not null;"`
	AccountUUID string `gorm:"not null;"`
	LastSeenAt  time.Time
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

type ChatRoomMemberInfo struct {
	ID                  uint
	ChatRoomID          uint
	AccountUUID         string
	LastSeenAt          time.Time
	Theme               string
	EnabledNotification bool
}

func (ChatRoomMember) Testing(accountUUID string, list *[]ChatRoomMember) *gorm.DB {
	return db.Joins(
		"JOIN chat_room_members ON chat_room_members.id = chat_room_member_settings.chat_room_member_id AND chat_room_members.account_uuid = ?", accountUUID,
	).Find(&list)
	// return db.Model(&ChatRoomInfo{}).Select(
	// 	"chat_room_members.id, chat_room_members.chat_room_id, chat_room_members.account_id, chat_room_members.last_seen_at",
	// ).Joins("left join chat_room_members on chat_room_members.id = chat_room_member_settings.chat_room_member_id").Scan(&c)
}
