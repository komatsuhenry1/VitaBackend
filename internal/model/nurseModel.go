package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nurse struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name" binding:"required"`
	Email            string             `bson:"email" json:"email" binding:"required,email"`
	Phone            string             `bson:"phone" json:"phone" binding:"required,phone"`
	Address          string             `bson:"address" json:"address" binding:"required,address"`
	Cpf              string             `bson:"cpf" json:"cpf" binding:"required"`
	PixKey           string             `bson:"pix_key" json:"pix_key" binding:"required"`
	Password         string             `bson:"password" json:"password" binding:"required"`
	VerificationSeal bool               `bson:"verification_seal" json:"verification_seal" binding:"required"`

	LicenseNumber   string `bson:"license_number" json:"license_number" binding:"required"` // registro profissional
	Specialization  string `bson:"specialization" json:"specialization"`                    // área (ex: pediatrics, geriatrics, ER)
	Shift           string `bson:"shift" json:"shift"`                                      // manhã, tarde, noite
	Department      string `bson:"department" json:"department"`                            // setor/hospital onde trabalha
	YearsExperience int    `bson:"years_experience" json:"years_experience"`

	Rating float64 `bson:"rating" json:"rating"`
	Price  float64 `bson:"price" json:"price"`
	Bio    string  `bson:"bio" json:"bio"`

	LicenseDocumentID     primitive.ObjectID `bson:"license_document_id,omitempty"`
	QualificationsID      primitive.ObjectID `bson:"qualifications_id,omitempty"`
	GeneralRegisterID     primitive.ObjectID `bson:"general_register_id,omitempty"`
	ResidenceComprovantId primitive.ObjectID `bson:"residence_comprovante_id,omitempty"`
	FaceImageID           primitive.ObjectID `bson:"face_image_id,omitempty"`

	Hidden      bool      `bson:"hidden" json:"hidden"`
	Role        string    `bson:"role" json:"role" binding:"required"`
	Online      bool      `bson:"online" json:"online" binding:"required"`
	FirstAccess bool      `bson:"first_access" json:"first_access"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	TempCode    int       `bson:"temp_code" json:"temp_code"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
