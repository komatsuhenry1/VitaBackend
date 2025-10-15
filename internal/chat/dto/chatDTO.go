package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ConversationDTO struct {
	PartnerID            primitive.ObjectID `bson:"partner_id" json:"partner_id"`
	PartnerName          string             `bson:"partner_name" json:"partner_name"`
	PartnerImageID       primitive.ObjectID `bson:"partner_image_id,omitempty" json:"partner_image_id,omitempty"`
	LastMessage          string             `bson:"last_message" json:"last_message"`
	LastMessageTimestamp time.Time          `bson:"last_message_timestamp" json:"last_message_timestamp"`
}