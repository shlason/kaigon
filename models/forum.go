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
	PopularTopics []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
