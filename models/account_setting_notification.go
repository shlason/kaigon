package models

type AccountSettingNotification struct {
	AccountID               uint
	AccountUUID             string
	FollowOrOwnArticleReply uint
	CommentTagged           bool
	ArticleTweet            bool
	CommentTweet            bool
	InterestRecommendation  bool
	Chat                    bool
	Followed                bool
}
