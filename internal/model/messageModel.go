package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderID   primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	ReceiverID primitive.ObjectID `bson:"receiver_id" json:"receiver_id"`
	Content    string             `bson:"content" json:"message"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Read       bool               `bson:"read" json:"read"`
}
