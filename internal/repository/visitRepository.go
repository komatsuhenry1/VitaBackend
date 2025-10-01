package repository

import(
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"medassist/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
	"fmt"
	"time"
)

type VisitRepository interface{
	CreateVisit(visit model.Visit) error
	FindAllVisitsForPatient(patientId string) ([]model.Visit, error)
	FindAllVisitsForNurse(nurseId string) ([]model.Visit, error)
	FindVisitById(id string) (model.Visit, error)
	UpdateVisitFields(id string, updates map[string]interface{}) (model.Visit,error)
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

func (r *visitRepository) FindAllVisitsForNurse(nurseId string) ([]model.Visit, error) {
	cursor, err := r.collection.Find(r.ctx, bson.M{"nurse_id": nurseId})
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

func (r *visitRepository) FindAllVisitsForPatient(patientId string) ([]model.Visit, error) {
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

func (r *visitRepository) FindVisitById(id string) (model.Visit, error) {
	var visit model.Visit

	// converter para ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return visit, fmt.Errorf("ID inválido: %w", err)
	}

	err = r.collection.FindOne(r.ctx, bson.M{"_id": objectID}).Decode(&visit)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Visit{}, fmt.Errorf("visita não encontrado")
		}
		return model.Visit{}, err
	}

	return visit, nil
}

func (r *visitRepository) UpdateVisitFields(id string, updates map[string]interface{}) (model.Visit,error) {
	cleanUpdates := bson.M{}

	for key, value := range updates {
		if value != nil || value == "" {
			cleanUpdates[key] = value
		}
	}

	if len(cleanUpdates) == 0 {
		return model.Visit{}, fmt.Errorf("nenhum campo válido para atualizar")
	}

	cleanUpdates["updated_at"] = time.Now()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Visit{}, fmt.Errorf("ID inválido")
	}

	update := bson.M{"$set": cleanUpdates}

	_, err = r.collection.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return model.Visit{}, err
	}

	return r.FindVisitById(id)
}