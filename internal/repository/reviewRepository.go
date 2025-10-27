package repository

import (
	"context"
	"medassist/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	CreateReview(review model.Review) error
}

type reviewRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewReviewRepository(db *mongo.Database) ReviewRepository {
	return &reviewRepository{
		collection: db.Collection("reviews"),
		ctx:        context.Background(),
	}
}

func (r *reviewRepository) CreateReview(review model.Review) error {
	_, err := r.collection.InsertOne(r.ctx, review)
	return err
}