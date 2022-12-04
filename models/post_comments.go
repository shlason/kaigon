package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const postCommentsCollectionName string = "post_comments"

type PostComment struct {
	MongoDBModel `bson:",inline"`
	PostID       primitive.ObjectID `bson:"post_id"`
	AccountID    uint               `bson:"account_id"`
	LikeCount    int                `bson:"like_count"`
	Content      string
}
