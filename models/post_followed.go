package models

const postFollowedCollectionName string = "post_followed"

type PostFollowed struct {
	MongoDBModel `bson:",inline"`
	PostID       uint `bson:"post_id"`
	AccountID    uint `bson:"account_id"`
}
