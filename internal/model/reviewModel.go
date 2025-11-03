package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReviewType string             `bson:"review_type" json:"review_type"`

	VisitId primitive.ObjectID `bson:"visit_id" json:"visit_id"`

	NurseName string             `bson:"nurse_name" json:"nurse_name"`
	NurseId   primitive.ObjectID `bson:"nurse_id" json:"nurse_id"`

	PatientName string             `bson:"patient_name" json:"patient_name"`
	PatientId   primitive.ObjectID `bson:"patient_id" json:"patient_id"`

	Rating  int    `bson:"rating" json:"rating"`
	Comment string `bson:"comment,omitempty" json:"comment,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
