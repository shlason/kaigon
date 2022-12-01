package models

const postFollowedCollectionName string = "post_followed"

type PostFollowed struct {
	mongoDBModel
	PostID    uint `bson:"post_id"`
	AccountID uint `bson:"account_id"`
}
