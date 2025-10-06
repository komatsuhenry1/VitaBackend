package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name" binding:"required"`
	Email          string             `bson:"email" json:"email" binding:"required,email"`
	Phone          string             `bson:"phone" json:"phone" binding:"required,phone"`
	Address        string             `bson:"address" json:"address" binding:"required,address"`
	Cpf            string             `bson:"cpf" json:"cpf" binding:"required"`
	Password       string             `bson:"password" json:"password" binding:"required"`
	Hidden         bool               `bson:"hidden" json:"hidden"`
	Role           string             `bson:"role" json:"role" binding:"required"`
	ProfileImageID primitive.ObjectID `bson:"profile_image_id" json:"profile_image_id"`
	FirstAccess    bool               `bson:"first_access" json:"first_access"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	TempCode       int                `bson:"temp_code" json:"temp_code"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
