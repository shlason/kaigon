package models

const postsCollectionName string = "posts"

type Post struct {
	mongoDBModel
	AccountID     uint
	IsAnonymous   bool
	ForumID       uint
	Title         string
	Content       string
	CommentCount  int
	ReactionStats map[string]int
}
