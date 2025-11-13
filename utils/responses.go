package utils

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
    adminDTO "medassist/internal/admin/dto"
    chatDTO "medassist/internal/chat/dto"
    nurseDTO "medassist/internal/nurse/dto"


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

// STRUCTS PARA O SWAGGER (possível passar para outro arquivo para manter organização)

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

type SuccessDashboardResponse struct {
    Success bool                         `json:"success" example:"true"`
    Message string                       `json:"message" example:"Dados de dashboard carregados com sucesso"`
    Data    adminDTO.DashboardAdminDataResponse `json:"data"`
}

// SuccessDocumentsResponse é a struct de sucesso para a lista de documentos
type SuccessDocumentsResponse struct {
    Success bool                       `json:"success" example:"true"`
    Message string                     `json:"message" example:"Documentos retornados com sucesso"`
    Data    []adminDTO.DocumentInfoResponse `json:"data"` // Note que é um slice []
}

type SuccessResponseString struct {
    Success bool   `json:"success" example:"true"`
    Message string `json:"message" example:"Operação realizada com sucesso"`
    Data    string `json:"data" example:"Mensagem de dados"`
}

type SuccessUserListsResponse struct {
    Success bool                  `json:"success" example:"true"`
    Message string                `json:"message" example:"Listas de usuários retornadas com sucesso."`
    Data    adminDTO.UserListsResponse `json:"data"`
}

// SuccessUserTypeResponse é a struct de sucesso para a resposta de UserTypeResponse
type SuccessUserTypeResponse struct {
    Success bool                 `json:"success" example:"true"`
    Message string               `json:"message" example:"Usuário atualizado com sucesso."`
    Data    adminDTO.UserTypeResponse `json:"data"`
}

// SuccessVisitTypeResponse é a struct de sucesso para a resposta de VisitTypeResponse
type SuccessVisitTypeResponse struct {
    Success bool                  `json:"success" example:"true"`
    Message string                `json:"message" example:"Visita atualizada com sucesso."`
    Data    adminDTO.VisitTypeResponse `json:"data"`
}

type SuccessMessagesResponse struct {
    Success bool             `json:"success" example:"true"`
    Message string           `json:"message" example:"Histórico de mensagens retornado com sucesso"`
    Data    []model.Message `json:"data"`
}

// SuccessConversationsResponse é a struct de sucesso para a lista de conversas
type SuccessConversationsResponse struct {
    Success bool                  `json:"success" example:"true"`
    Message string                `json:"message" example:"Lista de conversas retornada com sucesso"`
    Data    []chatDTO.ConversationDTO `json:"data"` // Reutilizável para Nurse e Patient
}

type SuccessNurseDashboardResponse struct {
    Success bool                              `json:"success" example:"true"`
    Message string                            `json:"message" example:"Dados de enfermeiro carregados com sucesso."`
    Data    nurseDTO.NurseDashboardDataResponseDTO `json:"data"`
}

type SuccessNurseResponse struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"Operação realizada com sucesso"`
    Data    model.Nurse `json:"data"`
}

// SuccessNurseVisitsListsResponse é a struct de sucesso para as listas de visitas
type SuccessNurseVisitsListsResponse struct {
    Success bool                    `json:"success" example:"true"`
    Message string                  `json:"message" example:"Visitas listadas com sucesso."`
    Data    nurseDTO.NurseVisitsListsDto `json:"data"`
}

type SuccessPatientProfileResponse struct {
    Success bool                          `json:"success" example:"true"`
    Message string                        `json:"message" example:"Perfil do paciente listado com sucesso."`
    Data    nurseDTO.PatientProfileResponseDTO `json:"data"`
}

// SuccessNurseUpdateResponse é a struct de sucesso para a atualização do enfermeiro
type SuccessNurseUpdateResponse struct {
    Success bool                         `json:"success" example:"true"`
    Message string                       `json:"message" example:"Usuário atualizado com sucesso."`
    Data    nurseDTO.NurseUpdateResponseDTO `json:"data"`
}

type SuccessAvailabilityResponse struct {
    Success bool                        `json:"success" example:"true"`
    Message string                      `json:"message" example:"Informações de disponibilidade carregadas com sucesso."`
    Data    nurseDTO.AvailabilityResponseDTO `json:"data"`
}

// SuccessNurseProfileResponse é a struct de sucesso para o perfil completo do enfermeiro
type SuccessNurseProfileResponse struct {
    Success bool                          `json:"success" example:"true"`
    Message string                        `json:"message" example:"Perfil completo de enfermeiro(a) listado com sucesso"`
    Data    nurseDTO.NurseProfileResponseDTO `json:"data"`
}

// SuccessNurseVisitInfoResponse é a struct de sucesso para os detalhes da visita
type SuccessNurseVisitInfoResponse struct {
    Success bool               `json:"success" example:"true"`
    Message string             `json:"message" example:"Informações de visita listadas com sucesso."`
    Data    nurseDTO.NurseVisitInfo `json:"data"`
}

type SuccessStripeOnboardingResponse struct {
    Success bool                            `json:"success" example:"true"`
    Message string                          `json:"message" example:"Link de onboarding criado com sucesso."`
    Data    nurseDTO.StripeOnboardingResponseDTO `json:"data"`
}