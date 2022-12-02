package models

const postsCollectionName string = "posts"

type Post struct {
	MongoDBModel  `bson:",inline"`
	AccountID     uint `bson:"account_id"`
	IsAnonymous   bool `bson:"is_anonymous"`
	ForumID       uint `bson:"forum_id"`
	Type          string
	Title         string
	Content       string
	Thumbnail     string
	MasterVision  string `bson:"master_vision"`
	CommentCount  int    `bson:"comment_count"`
	Topics        []string
	ReactionStats map[string]int `bson:"reaction_stats"`
}
