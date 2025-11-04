package repository

import (
	"context"
	"errors"
	"fmt"
	"medassist/internal/model"
	"medassist/utils"
	"strings"
	nurseDTO "medassist/internal/nurse/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitRepository interface {
	CreateVisit(visit model.Visit) error
	FindAllVisitsForPatient(patientId string) ([]model.Visit, error)
	FindAllVisitsForNurse(nurseId string) ([]model.Visit, error)
	FindAllPendingVisitsForNurse(nurseId string) ([]model.Visit, error)
	FindAllVisits() ([]model.Visit, error)
	FindVisitById(id string) (model.Visit, error)
	UpdateVisitFields(id string, updates map[string]interface{}) (model.Visit, error)
	DeleteVisit(visitId string) error
	FindAllCompletedVisitsForPatient(patientId string) ([]model.Visit, error)

	GetTotalVisitsCount() (int64, error)
    GetVisitsTodayCount() (int64, error)
    GetCompletedVisitsCountLast30Days() (int64, error)
    GetTotalRevenueLast30Days() (float64, error)
}

type visitRepository struct {
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

func (r *visitRepository) FindAllPendingVisitsForNurse(nurseId string) ([]model.Visit, error) {
	cursor, err := r.collection.Find(r.ctx, bson.M{"nurse_id": nurseId, "status": "PENDING"})
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

func (r *visitRepository) FindAllCompletedVisitsForPatient(patientId string) ([]model.Visit, error) {
	cursor, err := r.collection.Find(r.ctx, bson.M{"patient_id": patientId, "status": "COMPLETED"})
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

func (r *visitRepository) UpdateVisitFields(id string, updates map[string]interface{}) (model.Visit, error) {
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

func (r *visitRepository) FindAllVisits() ([]model.Visit, error) {
	var visits []model.Visit

	cursor, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		return visits, err
	}
	defer cursor.Close(r.ctx)

	if err = cursor.All(r.ctx, &visits); err != nil {
		return visits, err
	}

	return visits, nil
}

func (r *nurseRepository) UpdateVisitFields(id string, updates map[string]interface{}) (model.Visit, error) {
	var visit model.Visit

	fieldsToFormat := map[string]bool{
		"name": true,
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return visit, fmt.Errorf("ID inválido")
	}

	cleanUpdates := bson.M{}
	for key, value := range updates {
		if value == nil || value == "" {
			continue
		}

		valStr, ok := value.(string)
		if fieldsToFormat[strings.ToLower(key)] && ok {
			cleanUpdates[key] = utils.CapitalizeWords(valStr)
		} else {
			cleanUpdates[key] = value
		}
	}

	if len(cleanUpdates) == 0 {
		return visit, fmt.Errorf("nenhum campo válido para atualizar")
	}

	cleanUpdates["updated_at"] = time.Now()

	update := bson.M{"$set": cleanUpdates}

	_, err = r.collection.UpdateByID(r.ctx, objID, update)
	if err != nil {
		return visit, err
	}

	err = r.collection.FindOne(r.ctx, bson.M{"_id": objID}).Decode(&visit)
	return visit, err
}

func (r *visitRepository) DeleteVisit(visitId string) error {
	objID, err := primitive.ObjectIDFromHex(visitId)
	if err != nil {
		return fmt.Errorf("ID inválido")
	}

	_, err = r.collection.DeleteOne(r.ctx, bson.M{"_id": objID})
	return err
}

func (r *visitRepository) GetTotalVisitsCount() (int64, error) {
    return r.collection.CountDocuments(r.ctx, bson.M{})
}

// GetVisitsTodayCount retorna o número de visitas agendadas para hoje.
func (r *visitRepository) GetVisitsTodayCount() (int64, error) {
    // !! IMPORTANTE !!
    // Estou assumindo que seu model 'Visit' tem um campo 'visit_date' (time.Time)
    // que armazena a data do atendimento.
    // Se você usar 'created_at', contará os *agendamentos feitos hoje*,
    // o que é diferente de *atendimentos que acontecerão hoje*.
    
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    endOfDay := startOfDay.AddDate(0, 0, 1)

    // Troque "visit_date" pelo campo correto se for diferente
    filter := bson.M{"visit_date": bson.M{"$gte": startOfDay, "$lt": endOfDay}}
    
    // Se você *realmente* quiser os criados hoje, use:
    // filter := bson.M{"created_at": bson.M{"$gte": startOfDay, "$lt": endOfDay}}

    return r.collection.CountDocuments(r.ctx, filter)
}

// GetCompletedVisitsCountLast30Days retorna visitas concluídas nos últimos 30 dias.
func (r *visitRepository) GetCompletedVisitsCountLast30Days() (int64, error) {
    thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
    
    // Assumindo que quando uma visita é 'COMPLETED', o 'updated_at' é atualizado.
    // Se você tiver um campo 'completed_at', seria melhor usá-lo.
    filter := bson.M{
        "status":     "COMPLETED",
        "updated_at": bson.M{"$gte": thirtyDaysAgo},
    }
    return r.collection.CountDocuments(r.ctx, filter)
}

// GetTotalRevenueLast30Days usa aggregation para somar a receita de visitas concluídas.
func (r *visitRepository) GetTotalRevenueLast30Days() (float64, error) {
    // !! IMPORTANTE !!
    // Assumindo que seu model 'Visit' tem um campo 'price' (float64 ou similar)
    
    thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
    
    pipeline := mongo.Pipeline{
        bson.D{{"$match", bson.M{
            "status": "COMPLETED",
            // Novamente, assumindo 'updated_at' para data de conclusão
            "updated_at": bson.M{"$gte": thirtyDaysAgo}, 
        }}},
        bson.D{{"$group", bson.M{
            "_id":          nil,
            "totalRevenue": bson.M{"$sum": "$price"}, // Assumindo campo 'price'
        }}},
    }

    cursor, err := r.collection.Aggregate(r.ctx, pipeline)
    if err != nil {
        return 0, err
    }
    defer cursor.Close(r.ctx)

    if cursor.Next(r.ctx) {
        var result nurseDTO.TotalRevenueResult
        if err := cursor.Decode(&result); err != nil {
            return 0, err
        }
        return result.TotalRevenue, nil
    }

    return 0, nil // Retorna 0 se não houver receita
}
