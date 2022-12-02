package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const chatMessagesCollectionName string = "chat_messages"

type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	From      string
	To        uint
	Type      string
	Content   string
	Timestamp time.Time
}

func (m *ChatMessage) InsertOne() (*mongo.InsertOneResult, error) {
	return mdb.ChatMessages.InsertOne(context.TODO(), m)
}

func (ChatMessage) FindByTo(to uint) ([]*ChatMessage, error) {
	cursor, err := mdb.ChatMessages.Find(context.TODO(), bson.D{{Key: "to", Value: to}}, options.Find())

	if err != nil {
		return []*ChatMessage{}, err
	}

	var results []*ChatMessage

	for cursor.Next(context.TODO()) {
		var elem ChatMessage
		err := cursor.Decode(&elem)
		if err != nil {
			return []*ChatMessage{}, err
		}

		results = append(results, &elem)
	}

	return results, nil
}
