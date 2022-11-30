package models

const forumsCollectionName string = "forums"

type Forum struct {
	mongoDBModel
	Name          string
	Icon          string
	Banner        string
	Rule          string
	Description   string
	PopularTopics []string `bson:"popular_topics"`
}
