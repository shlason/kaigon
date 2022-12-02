package models

const topicsCollectionName string = "topics"

type Topic struct {
	MongoDBModel  `bson:",inline"`
	Name          string
	PostCount     int `bson:"post_count"`
	FollowedCount int `bson:"followed_count"`
}
