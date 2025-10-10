package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthUser struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string             `bson:"name"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
	TwoFactor        bool               `bson:"two_factor"`
	VerificationSeal bool               `bson:"verification_seal"`
	Role             string             `bson:"role"`
	Hidden           bool               `bson:"hidden"`
	TempCode         int                `bson:"temp_code"`
	ProfileImageID   primitive.ObjectID `bson:"profile_image_id"`
}

type CodeResponseDTO struct {
	Code int `json:"code"`
}
type FileData struct {
	Data        []byte
	ContentType string
	Filename    string
}