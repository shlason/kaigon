package models

import "gorm.io/gorm"

type AccountSettingNotification struct {
	AccountID               uint   `gorm:"unique; not null;"`
	AccountUUID             string `gorm:"unique; not null;"`
	FollowOrOwnArticleReply uint
	CommentTagged           bool
	ArticleTweet            bool
	CommentTweet            bool
	InterestRecommendation  bool
	Chat                    bool
	Followed                bool
}

func (accountSettingNotification *AccountSettingNotification) Create() *gorm.DB {
	return db.Create(&accountSettingNotification)
}
