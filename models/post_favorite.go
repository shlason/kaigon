package models

const postFavoriteCollectionName string = "post_favorite"

type PostFavorite struct {
	mongoDBModel
	AccountID uint `bson:"account_id"`
	PostID    uint `bson:"post_id"`
}
