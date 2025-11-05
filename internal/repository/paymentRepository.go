package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
}

type paymentRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewPaymentRepository(db *mongo.Database) PaymentRepository {
	return &reviewRepository{
		collection: db.Collection("payments"),
		ctx:        context.Background(),
	}
}
