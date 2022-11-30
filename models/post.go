package models

const postsCollectionName string = "posts"

type Post struct {
	mongoDBModel
	AccountID     uint `bson:"account_id"`
	IsAnonymous   bool `bson:"is_anonymous"`
	ForumID       uint `bson:"forum_id"`
	Title         string
	Content       string
	CommentCount  int `bson:"comment_count"`
	Topics        []string
	ReactionStats map[string]int `bson:"reaction_stats"`
}
