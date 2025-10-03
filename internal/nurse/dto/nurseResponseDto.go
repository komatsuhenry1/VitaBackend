package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NurseVisitsListsDto struct {
	Pending   []VisitDto `json:"pending"`
	Confirmed []VisitDto `json:"confirmed"`
	Completed []VisitDto `json:"completed"`
}

type VisitDto struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Reason      string  `json:"reason"`
	VisitType   string  `json:"visit_type"`
	VisitValue  float64 `json:"visit_value"`
	CreatedAt   string  `json:"created_at"`
	Date        string  `json:"date"`
	Status      string  `json:"status"`
	PatientName string  `json:"patient_name"`
	PatientId   string  `json:"patient_id"`
	NurseName   string  `json:"nurse_name"`
}

type PatientProfileResponseDTO struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	Address     string             `json:"address"`
	Cpf         string             `json:"cpf"`
	Password    string             `json:"password"`
	Hidden      bool               `json:"hidden"`
	Role        string             `json:"role"`
	FirstAccess bool               `json:"first_access"`
	CreatedAt   time.Time          `json:"created_at"`
	TempCode    int                `json:"temp_code"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
