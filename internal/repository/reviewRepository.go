package repository

import (
	"context"
	"medassist/internal/model"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	CreateReview(review model.Review) error
	FindReviewByVisitId(visitId string) (model.Review, error)
	FindAverageRatingByNurseId(nurseId string) (float64, error)
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

func (r *reviewRepository) FindAverageRatingByNurseId(nurseId string) (float64, error) {

	objID, err := primitive.ObjectIDFromHex(nurseId)
	if err != nil {
		fmt.Printf("ID do enfermeiro(a) inválido: %s, erro: %v\n", nurseId, err)
		return 0.0, nil
	}

	matchStage := bson.D{{Key: "$match", Value: bson.M{"nurse_id": objID}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":           nil,
			"averageRating": bson.M{"$avg": "$rating"}, // <- Nome é "averageRating"
		}},
	}

	pipeline := mongo.Pipeline{matchStage, groupStage}

	cursor, err := r.collection.Aggregate(r.ctx, pipeline)
	if err != nil {
		fmt.Printf("Erro ao agregar ratings: %v\n", err)
		return 0.0, err
	}
	defer cursor.Close(r.ctx)

	// MUDANÇA AQUI
	var result struct {
		// O nome na tag 'bson' deve ser idêntico ao nome no $group
		AverageRating float64 `bson:"averageRating"`
	}
	// FIM DA MUDANÇA

	if cursor.Next(r.ctx) {
		if err = cursor.Decode(&result); err != nil {
			fmt.Printf("Erro ao decodificar média de rating: %v\n", err)
			return 0.0, err
		}
	}

	return result.AverageRating, nil
}
