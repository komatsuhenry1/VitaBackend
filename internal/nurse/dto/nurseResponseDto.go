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
	ID             string  `json:"id"`
	Description    string  `json:"description"`
	Reason         string  `json:"reason"`
	VisitType      string  `json:"visit_type"`
	VisitValue     float64 `json:"visit_value"`
	CreatedAt      string  `json:"created_at"`
	Date           string  `json:"date"`
	Status         string  `json:"status"`
	PatientName    string  `json:"patient_name"`
	PatientId      string  `json:"patient_id"`
	NurseName      string  `json:"nurse_name"`
	PatientImageID string  `json:"patient_image_id"`
}

type PatientProfileResponseDTO struct {
	ID             primitive.ObjectID `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Phone          string             `json:"phone"`
	Address        string             `json:"address"`
	Cpf            string             `json:"cpf"`
	Password       string             `json:"password"`
	Hidden         bool               `json:"hidden"`
	Role           string             `json:"role"`
	FirstAccess    bool               `json:"first_access"`
	CreatedAt      time.Time          `json:"created_at"`
	TempCode       int                `json:"temp_code"`
	UpdatedAt      time.Time          `json:"updated_at"`
	ProfileImageID string             `json:"profile_image_id"`
}

type StatsDTO struct {
	PatientsAttended  int     `json:"patients_attended"`
	AppointmentsToday int     `json:"appointments_today"`
	AverageRating     float64 `json:"average_rating"`
	MonthlyEarnings   float64 `json:"monthly_earnings"`
}

type NurseProfileResponseDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Coren           string `json:"coren"`
	ExperienceYears int    `json:"experience_years"`
	Department      string `json:"department"`
	Bio             string `json:"bio"`
}

type AvailabilityDTO struct {
	IsAvailable    bool   `json:"is_available"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Specialization string `json:"specialization"`
}

type NurseDashboardDataResponseDTO struct {
	Online bool       `json:"online"`
	Stats  StatsDTO   `json:"stats"`
	Visits []VisitDto `json:"visits"`
	// History      []VisitDto                  `json:"history"`
	Profile      NurseProfileResponseDTO `json:"profile"`
	Availability AvailabilityDTO         `json:"availability"`
}

type NurseUpdateResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `bson:"email" json:"email"`
	Hidden    bool      `bson:"hidden" json:"hidden"`
	Role      string    `bson:"role" json:"role"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type AvailabilityResponseDTO struct {
	Online                 bool     `json:"online"`
	StartTime              string   `json:"start_time"`
	EndTime                string   `json:"end_time"`
	Specialization         string   `json:"specialization"`
	Price                  float64  `json:"price"`
	MaxPatientsPerDay      int      `json:"max_patients_per_day"`
	DaysAvailable          []string `json:"days_available"`
	Services               []string `json:"services"`
	AvailableNeighborhoods []string `json:"available_neighborhoods"`
}
