package models

const postFollowedCollectionName string = "post_followed"

type PostFollowed struct {
	mongoDBModel
	AccountID uint `bson:"account_id"`
}
