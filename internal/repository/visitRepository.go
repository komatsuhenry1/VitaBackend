package repository

import(
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"medassist/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type VisitRepository interface{
	CreateVisit(visit model.Visit) error
	FindAllVisits(patientId string) ([]model.Visit, error)
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

func (r *visitRepository) FindAllVisits(patientId string) ([]model.Visit, error) {
	cursor, err := r.collection.Find(r.ctx, bson.M{"patient_id": patientId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var visits []model.Visit
	if err := cursor.All(context.TODO(), &visits); err != nil {
		return nil, err
	}
	return visits, nil
}