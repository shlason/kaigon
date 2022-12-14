package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

// Redis
type ChatRoomInviteCode struct {
	Code       string
	ChatRoomID uint
}

func (cric *ChatRoomInviteCode) Create() error {
	cric.Code = uuid.NewString()
	return rdb.SetNX(rctx, fmt.Sprintf("chat:room:invite:code:%s", cric.Code), cric.ChatRoomID, 10*time.Minute).Err()
}

func (cric *ChatRoomInviteCode) Read() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("chat:room:invite:code:%s", cric.Code)).Result()

	if err != nil {
		return err
	}

	chatRoomID, err := strconv.ParseUint(val, 10, 22)

	if err != nil {
		return err
	}

	cric.ChatRoomID = uint(chatRoomID)
	return nil
}
