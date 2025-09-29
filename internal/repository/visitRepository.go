package repository

import(
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"medassist/internal/model"
)

type VisitRepository interface{
	CreateVisit(visit model.Visit) error
}

type visitRepository struct{
	collection *mongo.Collection
	ctx        context.Context

}

func NewVisitRepository(db *mongo.Database) VisitRepository {
	return &visitRepository{
		collection: db.Collection("visits"),
		ctx:        context.Background(),
	}
}

func (r *visitRepository) CreateVisit(visit model.Visit) error {
	_, err := r.collection.InsertOne(r.ctx, visit)
	return err
}