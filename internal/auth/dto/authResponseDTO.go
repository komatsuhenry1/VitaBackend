package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthUser struct {
	ID               primitive.ObjectID `bson:"id" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
	TwoFactor        bool               `bson:"two_factor" json:"two_factor"`
	VerificationSeal bool               `bson:"verification_seal" json:"verification_seal"`
	Role             string             `bson:"role" json:"role"`
	Hidden           bool               `bson:"hidden" json:"hidden"`
	TempCode         int                `bson:"temp_code" json:"temp_code"`
	ProfileImageID   primitive.ObjectID `bson:"profile_image_id" json:"profile_image_id"`
}

type CodeResponseDTO struct {
	Code int `json:"code"`
}
type FileData struct {
	Data        []byte
	ContentType string
	Filename    string
}