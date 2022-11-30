package models

const topicFollowedCollectionName string = "topic_followed"

type TopicFollowed struct {
	mongoDBModel
	AccountID uint `bson:"account_id"`
}
