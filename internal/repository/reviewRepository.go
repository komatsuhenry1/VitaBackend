package repository

import (
	"context"
	"medassist/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	CreateReview(review model.Review) error
	FindReviewByVisitId(visitId string) (model.Review, error)
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

func (r *reviewRepository) FindReviewByVisitId(visitId string) (model.Review, error) {
	var review model.Review

	// MUDANÇA: Converter a string 'visitId' de volta para um ObjectID
	objID, err := primitive.ObjectIDFromHex(visitId)
	if err != nil {
		// Se a string do ID for inválida (ex: "123"), ela não existe no banco.
		// Retornar ErrNoDocuments é seguro, pois o service já sabe tratar.
		return model.Review{}, mongo.ErrNoDocuments
	}

	// MUDANÇA: Usar o 'objID' (do tipo ObjectID) na consulta, em vez da string 'visitId'
	err = r.collection.FindOne(r.ctx, bson.M{"visit_id": objID}).Decode(&review)
	if err != nil {
		// O service 'FindAllVisits' vai tratar o 'mongo.ErrNoDocuments' se não achar
		return model.Review{}, err
	}

	return review, nil
}
