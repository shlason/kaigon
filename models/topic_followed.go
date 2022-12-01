package models

const topicFollowedCollectionName string = "topic_followed"

type TopicFollowed struct {
	mongoDBModel
	TopicID   uint `bson:"topic_id"`
	AccountID uint `bson:"account_id"`
}
