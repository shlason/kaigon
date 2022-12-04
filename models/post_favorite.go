package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const postFavoriteCollectionName string = "post_favorite"

type PostFavorite struct {
	MongoDBModel `bson:",inline"`
	PostID       primitive.ObjectID `bson:"post_id"`
	AccountID    uint               `bson:"account_id"`
}
