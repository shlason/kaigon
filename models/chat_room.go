package models

import (
	"strings"

	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	Type             string `gorm:"not null;"`
	MaximumMemberNum int    `gorm:"default:50;not null;"`
	Emoji            string
	Name             string
	Avatar           string
}

func (cr *ChatRoom) Create() *gorm.DB {
	return db.Create(&cr)
}

func (ChatRoom) ReadAllByIDs(ids []interface{}, list *[]ChatRoom) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "id = ?")
	}

	return db.Where(strings.Join(fields, " OR "), ids...).Find(&list)
}

func (cr *ChatRoom) UpdateByID(m map[string]interface{}) *gorm.DB {
	return db.Model(&cr).Where("id = ?", cr.ID).Updates(m)
}
