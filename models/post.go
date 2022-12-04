package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const postsCollectionName string = "posts"

type Post struct {
	MongoDBModel  `bson:",inline"`
	AccountID     uint               `bson:"account_id"`
	IsAnonymous   bool               `bson:"is_anonymous"`
	ForumID       primitive.ObjectID `bson:"forum_id"`
	Type          string
	Title         string
	Content       string
	Thumbnail     string
	MasterVision  string `bson:"master_vision"`
	CommentCount  int    `bson:"comment_count"`
	Topics        []string
	ReactionStats map[string]int `bson:"reaction_stats"`
}

func (p *Post) InsertOne() (*mongo.InsertOneResult, error) {
	current := time.Now().UTC()
	p.CreatedAt = current
	p.UpdatedAt = current
	return mdb.Posts.InsertOne(context.TODO(), p)
}
