package models

const postFavoriteCollectionName string = "post_favorite"

type PostFavorite struct {
	MongoDBModel `bson:",inline"`
	PostID       uint `bson:"post_id"`
	AccountID    uint `bson:"account_id"`
}
