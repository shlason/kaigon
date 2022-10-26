package account

import "github.com/shlason/kaigon/controllers"

type getSettingNotificationResponsePayload struct {
	FollowOrOwnArticleReply uint `json:"followOrOwnArticleReply"`
	CommentTagged           bool `json:"commentTagged"`
	ArticleTweet            bool `json:"articleTweet"`
	CommentTweet            bool `json:"commentTweet"`
	InterestRecommendation  bool `json:"interestRecommendation"`
	Chat                    bool `json:"chat"`
	Followed                bool `json:"followed"`
}

type putSettingNotificationRequestPayload struct {
	FollowOrOwnArticleReply uint `json:"followOrOwnArticleReply"`
	CommentTagged           bool `json:"commentTagged"`
	ArticleTweet            bool `json:"articleTweet"`
	CommentTweet            bool `json:"commentTweet"`
	InterestRecommendation  bool `json:"interestRecommendation"`
	Chat                    bool `json:"chat"`
	Followed                bool `json:"followed"`
}

func (p *putSettingNotificationRequestPayload) check() (errResponse controllers.JSONResponse, inNotValid bool) {
	var acceptFollowOrOwnArticleReply = map[uint]string{
		0: "Close",
		1: "Open",
		2: "Only Tagged",
	}

	if _, ok := acceptFollowOrOwnArticleReply[p.FollowOrOwnArticleReply]; !ok {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}
