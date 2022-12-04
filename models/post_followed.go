package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const postFollowedCollectionName string = "post_followed"

type PostFollowed struct {
	MongoDBModel `bson:",inline"`
	PostID       primitive.ObjectID `bson:"post_id"`
	AccountID    uint               `bson:"account_id"`
}
