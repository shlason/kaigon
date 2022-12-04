package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const topicFollowedCollectionName string = "topic_followed"

type TopicFollowed struct {
	MongoDBModel `bson:",inline"`
	TopicID      primitive.ObjectID `bson:"topic_id"`
	AccountID    uint               `bson:"account_id"`
}
