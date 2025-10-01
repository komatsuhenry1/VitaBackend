package dto

import (
	"time"
)

type ContactUsDTO struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type CreateVisitDto struct {
	Description string `json:"description" binding:"required"`
	Reason      string `json:"reason" binding:"required"`

	NurseId string `json:"nurse_id" binding:"required"`

	VisitValue float64   `json:"value" binding:"required"`
	VisitType  string    `json:"visit_type" binding:"required"`
	VisitDate  time.Time `json:"date" binding:"required"`
}
