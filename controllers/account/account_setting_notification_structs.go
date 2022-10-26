package account

type getSettingNotificationResponsePayload struct {
	FollowOrOwnArticleReply uint `json:"followOrOwnArticleReply"`
	CommentTagged           bool `json:"commentTagged"`
	ArticleTweet            bool `json:"articleTweet"`
	CommentTweet            bool `json:"commentTweet"`
	InterestRecommendation  bool `json:"interestRecommendation"`
	Chat                    bool `json:"chat"`
	Followed                bool `json:"followed"`
}
