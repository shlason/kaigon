package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (Forum) Find() ([]Forum, error) {
	cursor, err := mdb.Forums.Find(context.TODO(), bson.D{}, options.Find())

	if err != nil {
		return []Forum{}, err
	}

	var results []Forum

	for cursor.Next(context.TODO()) {
		var elem Forum
		err := cursor.Decode(&elem)
		fmt.Println(elem)
		if err != nil {
			return []Forum{}, err
		}

		results = append(results, elem)
	}

	return results, nil
}

func (f *Forum) FindOneByName() error {
	return mdb.Forums.FindOne(context.TODO(), bson.D{{Key: "name", Value: f.Name}}).Decode(&f)
}

func (f *Forum) FindOneByID() error {
	return mdb.Forums.FindOne(context.TODO(), bson.D{{Key: "_id", Value: f.ID}}).Decode(&f)
}
