package utils

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"


	model "medassist/internal/model"
)

func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SendErrorResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
	})
}

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param : %s (type: %s) is required", name, typ)
}

// --- 💡 ADICIONE ESTAS STRUCTS PARA O SWAGGER ---

// ErrorResponse define a estrutura padrão para respostas de erro da API.
// O Swag usará isso para documentar o @Failure.
type ErrorResponse struct {
    Success bool   `json:"success" example:"false"`
    Message string `json:"message" example:"Descrição do erro"`
}

// SuccessResponseUser define a estrutura de sucesso para o endpoint de registro de usuário.
// O Swag usará isso para documentar o @Success.
type SuccessResponseUser struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Usuário criado com sucesso"`
    Data    model.User `json:"data"` // Aponta para a struct real do usuário
}

type SuccessResponseNoData struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Operação realizada com sucesso"`
}

type SuccessResponseNurse struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Cadastro de enfermeiro solicitado com sucesso."`
    Data    model.Nurse `json:"data"` // Aponta para a struct real do enfermeiro
}

// AuthUserResponse é a struct parcial do usuário para login/validação
type AuthUserResponse struct {
    ID             primitive.ObjectID `json:"_id"`
    Name           string             `json:"name"`
    Email          string             `json:"email"`
    Role           string             `json:"role"`
    TwoFactor      bool               `json:"two_factor"`
    ProfileImageID primitive.ObjectID `json:"profile_image_id"`
}

// ValidateCodeResponse é a struct de dados para validação de código
type ValidateCodeResponse struct {
    Token string           `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User  AuthUserResponse `json:"user"`
}

// SuccessValidateCodeResponse é a struct de sucesso para validação de código
type SuccessValidateCodeResponse struct {
    Success bool                 `json:"success" example:"true"`
    Message string               `json:"message" example:"Código validado com sucesso"`
    Data    ValidateCodeResponse `json:"data"`
}