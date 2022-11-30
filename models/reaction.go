package models

const reactionsCollectionName string = "reactions"

type Reaction struct {
	mongoDBModel
	AccountID uint `bson:"account_id"`
	PostID    uint `bson:"post_id"`
	Type      string
}
