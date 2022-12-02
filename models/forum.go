package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const forumsCollectionName string = "forums"

type Forum struct {
	MongoDBModel  `bson:",inline"`
	Name          string
	Icon          string
	Banner        string
	Rule          string
	Description   string
	PopularTopics []string `bson:"popular_topics"`
}

func (f *Forum) InsertOne() (*mongo.InsertOneResult, error) {
	current := time.Now().UTC()
	f.CreatedAt = current
	f.UpdatedAt = current
	return mdb.Forums.InsertOne(context.TODO(), f)
}

func (f *Forum) FindOneByName() error {
	return mdb.Forums.FindOne(context.TODO(), bson.D{{Key: "name", Value: f.Name}}).Decode(&f)
}
