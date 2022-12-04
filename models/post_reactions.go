package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const postReactionsCollectionName string = "post_reactions"

type PostReaction struct {
	MongoDBModel `bson:",inline"`
	AccountID    uint               `bson:"account_id"`
	PostID       primitive.ObjectID `bson:"post_id"`
	Type         string
}
