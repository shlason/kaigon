package models

import (
	"strings"

	"gorm.io/gorm"
)

type AccountSettingNotification struct {
	ID                      uint   `gorm:"primarykey"`
	AccountID               uint   `gorm:"unique; not null;"`
	AccountUUID             string `gorm:"unique; not null;"`
	FollowOrOwnArticleReply uint   `gorm:"default:1"`
	CommentTagged           bool   `gorm:"default:true"`
	ArticleTweet            bool   `gorm:"default:true"`
	CommentTweet            bool   `gorm:"default:true"`
	InterestRecommendation  bool   `gorm:"default:true"`
	Chat                    bool   `gorm:"default:true"`
	Followed                bool   `gorm:"default:true"`
}

func (accountSettingNotification *AccountSettingNotification) Create() *gorm.DB {
	return db.Create(&accountSettingNotification)
}

func (accountSettingNotification *AccountSettingNotification) ReadByAccountUUID() *gorm.DB {
	return db.First(&accountSettingNotification, "account_uuid = ?", accountSettingNotification)
}

func (accountSettingNotification *AccountSettingNotification) UpdateByAccountUUID(m map[string]interface{}) *gorm.DB {
	return db.Model(&accountSettingNotification).Where("account_uuid = ?", accountSettingNotification.AccountUUID).Updates(m)
}

func (AccountSettingNotification) DeleteByAccountIDs(ids []interface{}) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "account_id = ?")
	}

	return db.Where(strings.Join(fields, " OR "), ids...).Delete(&AccountSettingNotification{})
}
