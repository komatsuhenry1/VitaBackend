package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nurse struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                   string             `bson:"name" json:"name" binding:"required"`
	Email                  string             `bson:"email" json:"email" binding:"required,email"`
	Phone                  string             `bson:"phone" json:"phone" binding:"required,phone"`
	Cpf                    string             `bson:"cpf" json:"cpf" binding:"required"`
	PixKey                 string             `bson:"pix_key" json:"pix_key" binding:"required"`
	Password               string             `bson:"password" json:"password" binding:"required"`
	TwoFactor              int                `bson:"two_factor" json:"two_factor"`
	VerificationSeal       bool               `bson:"verification_seal" json:"verification_seal" binding:"required"`
	MaxPatientsPerDay      int                `bson:"max_patients_per_day" json:"max_patients_per_day"`
	DaysAvailable          []string           `bson:"days_available" json:"days_available"`
	Services               []string           `bson:"services" json:"services"`
	AvailableNeighborhoods []string           `bson:"available_neighborhoods" json:"available_neighborhoods"`

	Address      string `bson:"address" json:"address" binding:"required,address"`
	CEP          string `json:"cep"`
	Street       string `bson:"street" json:"street" binding:"required"`
	Number       string `bson:"number" json:"number" binding:"required"`
	Complement   string `bson:"complement" json:"complement"`
	Neighborhood string `bson:"neighborhood" json:"neighborhood" binding:"required"`
	City         string `bson:"city" json:"city" binding:"required"`
	UF           string `bson:"uf" json:"uf" binding:"required"`

	Coren           string `bson:"coren" json:"coren" binding:"required"` // registro profissional
	Specialization  string `bson:"specialization" json:"specialization"`  // área (ex: pediatrics, geriatrics, ER)
	Shift           string `bson:"shift" json:"shift"`                    // manhã, tarde, noite
	Department      string `bson:"department" json:"department"`          // setor/hospital onde trabalha
	YearsExperience int    `bson:"years_experience" json:"years_experience"`

	Rating float64 `bson:"rating" json:"rating"`
	Price  float64 `bson:"price" json:"price"`
	Bio    string  `bson:"bio" json:"bio"`

	LicenseDocumentID     primitive.ObjectID `bson:"license_document_id" json:"license_document_id" binding:"required"`
	QualificationsID      primitive.ObjectID `bson:"qualifications_id" json:"qualifications_id" binding:"required"`
	GeneralRegisterID     primitive.ObjectID `bson:"general_register_id" json:"general_register_id" binding:"required"`
	ResidenceComprovantId primitive.ObjectID `bson:"residence_comprovant_id" json:"residence_comprovant_id" binding:"required"`
	ProfileImageID        primitive.ObjectID `bson:"profile_image_id" json:"profile_image_id" binding:"required"`

	Hidden      bool      `bson:"hidden" json:"hidden"`
	Role        string    `bson:"role" json:"role" binding:"required"`
	Online      bool      `bson:"online" json:"online" binding:"required"`
	FirstAccess bool      `bson:"first_access" json:"first_access"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	StartTime   string    `bson:"start_time" json:"start_time" binding:"required"`
	EndTime     string    `bson:"end_time" json:"end_time" binding:"required"`
	TempCode    int       `bson:"temp_code" json:"temp_code"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
