package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (f *Forum) InsertOne() (*mongo.InsertOneResult, error) {
	return mdb.Forums.InsertOne(context.TODO(), f)
}

func (f *Forum) FindOneByName() error {
	return mdb.Forums.FindOne(context.TODO(), bson.D{{Key: "name", Value: f.Name}}).Decode(&f)
}
