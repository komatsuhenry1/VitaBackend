package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visit struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status           string             `bson:"status" json:"status" binding:"required"`
	ConfirmationCode string             `bson:"confirmation_code" json:"confirmation_code"`

	PatientId    string `bson:"patient_id" json:"patient_id" binding:"required"`
	PatientName  string `bson:"patient_name" json:"patient_name" binding:"required"`
	PatientEmail string `bson:"patient_email" json:"patient_email" binding:"required"`

	CEP          string `bson:"cep" json:"cep" binding:"required,address"`
	Street       string `bson:"street" json:"street" binding:"required"`
	Number       string `bson:"number" json:"number" binding:"required"`
	Complement   string `bson:"complement" json:"complement"`
	Neighborhood string `bson:"neighborhood" json:"neighborhood" binding:"required"`

	Description  string `bson:"description" json:"description" binding:"required"`
	Reason       string `bson:"reason" json:"reason" binding:"required"`
	CancelReason string `bson:"cancel_reason" json:"cancel_reason"`

	NurseId   string `bson:"nurse_id" json:"nurse_id" binding:"required"`
	NurseName string `bson:"nurse_name" json:"nurse_name" binding:"required"`

	VisitValue       float64   `bson:"value" json:"value" binding:"required"`
	VisitType        string    `bson:"visit_type" json:"visit_type" binding:"required"`
	VisitRequestType string    `bson:"visit_request_type" json:"visit_request_type" binding:"required"`
	VisitDate        time.Time `bson:"visit_date" json:"visit_date" binding:"required"`

	PaymentIntentID string `bson:"payment_intent_id" json:"payment_intent_id" binding:"required"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
