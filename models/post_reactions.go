package models

const postReactionsCollectionName string = "post_reactions"

type PostReaction struct {
	mongoDBModel
	AccountID uint `bson:"account_id"`
	PostID    uint `bson:"post_id"`
	Type      string
}
