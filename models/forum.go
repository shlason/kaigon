package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const forumsCollectionName string = "forums"

type Forum struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string
	Icon          string
	Banner        string
	Rule          string
	Description   string
	PopularTopics []string  `bson:"popular_topics"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
	DeletedAt     time.Time `bson:"deleted_at"`
}
