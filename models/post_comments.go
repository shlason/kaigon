package models

const postCommentsCollectionName string = "post_comments"

type PostComment struct {
	mongoDBModel
	PostID    uint `bson:"post_id"`
	AccountID uint `bson:"account_id"`
	LikeCount int  `bson:"like_count"`
	Content   string
}
