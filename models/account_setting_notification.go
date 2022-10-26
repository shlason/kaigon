package models

import "gorm.io/gorm"

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
