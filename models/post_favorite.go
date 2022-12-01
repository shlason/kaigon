package models

const postFavoriteCollectionName string = "post_favorite"

type PostFavorite struct {
	mongoDBModel
	PostID    uint `bson:"post_id"`
	AccountID uint `bson:"account_id"`
}
