package nurse

import (
	"fmt"
	"medassist/internal/nurse/dto"
	"medassist/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type NurseHandler struct {
	nurseService NurseService
}

func NewNurseHandler(nurseService NurseService) *NurseHandler {
	return &NurseHandler{nurseService: nurseService}
}

// @Summary Dashboard do Enfermeiro
// @Description Retorna os dados principais para o dashboard do enfermeiro logado (perfil, stats, visitas, etc.). Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessNurseDashboardResponse "Dados de dashboard carregados com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar dados"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/dashboard [get]
func (h *NurseHandler) NurseDashboard(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	dashboardData, err := h.nurseService.NurseDashboardData(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Dados de enfermeiro carregados com sucesso.", dashboardData)
}

// @Summary Perfil completo do Enfermeiro (para Dashboard)
// @Description Retorna o perfil completo do enfermeiro logado, incluindo reviews, dados privados, e estatísticas. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessNurseProfileResponse "Perfil completo de enfermeiro(a) listado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar perfil"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/dashboard_info [get]
func (h *NurseHandler) NurseDashboardData(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseProfile, err := h.nurseService.GetNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

// @Summary Altera o status online do Enfermeiro
// @Description Ativa ou desativa o status 'online' do enfermeiro logado (toggle). Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessNurseResponse "Serviço ativado com sucesso (retorna o objeto Nurse atualizado)"
// @Failure 400 {object} utils.ErrorResponse "Erro ao atualizar status"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/online [patch]
func (h *NurseHandler) ChangeOnlineNurse(c *gin.Context) {
	nurseId := utils.GetUserId(c) // pega pelo token

	//VALIDACAO DE ROLE
	claims, exists := c.Get("claims")
	if !exists {
		utils.SendErrorResponse(c, "Usuário não autenticado.", http.StatusUnauthorized)
		return
	}
	role, ok := claims.(jwt.MapClaims)["role"].(string)
	if !ok {
		utils.SendErrorResponse(c, "Usuário não autenticado.", http.StatusUnauthorized)
		return
	}

	if role != "NURSE" {
		utils.SendErrorResponse(c, "Rota apenas para usuários comuns.", http.StatusUnauthorized)
		return
	}

	nurseStatus, err := h.nurseService.UpdateAvailablityNursingService(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Serviço ativado com sucesso.", nurseStatus)
}

// @Summary Lista de visitas do Enfermeiro
// @Description Retorna todas as visitas (pendentes, confirmadas, concluídas, etc.) associadas ao enfermeiro logado. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessNurseVisitsListsResponse "Visitas listadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao buscar visitas"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/visits [get]
func (h *NurseHandler) GetAllVisits(c *gin.Context) {
	//pega o id pelo token
	nurseId := utils.GetUserId(c)

	visits, err := h.nurseService.GetAllVisits(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Visitas [PENDENTE/MARCADAS/CONCLUIDAS] listadas com sucesso.", visits)
}

// @Summary Confirma ou cancela uma visita (Enfermeiro)
// @Description Permite ao enfermeiro confirmar uma visita PENDENTE/REJEITADA, ou cancelar uma visita CONFIRMADA (com motivo). Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser modificada"
// @Param payload body dto.CancelReason false "Motivo do cancelamento. Obrigatório apenas ao cancelar uma visita 'CONFIRMED'. Pode ser um JSON vazio {} para confirmar."
// @Success 200 {object} utils.SuccessResponseString "Status da visita alterado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido, JSON inválido ou erro na lógica de status"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/visit/{id} [patch]
func (h *NurseHandler) ConfirmOrCancelVisit(c *gin.Context) {
	nurseId := utils.GetUserId(c) // pega o id do user pela req

	fmt.Println("nurseiD", nurseId)

	visitId := c.Param("id")
	// permite que o campo 'reason' venha vazio
	var reason dto.CancelReason
	if err := c.ShouldBindJSON(&reason); err != nil {
		if err.Error() != "EOF" {
			utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
	}

	response, err := h.nurseService.ConfirmOrCancelVisit(nurseId, visitId, reason.Reason)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}
	fmt.Println("response", response)

	utils.SendSuccessResponse(c, "Status da visita alterado com sucesso.", response)
}

// @Summary Perfil do Paciente
// @Description Retorna o perfil público de um paciente (usado por enfermeiros). Requer autenticação de Enfermeiro ou Paciente.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Paciente para ver o perfil"
// @Success 200 {object} utils.SuccessPatientProfileResponse "Perfil do paciente listado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao buscar perfil"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /nurse/patient/{id} [get]
func (h *NurseHandler) GetPatientProfile(c *gin.Context) {

	patientId := c.Param("id")

	patientProfile, err := h.nurseService.GetPatientProfile(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Perfil do paciente listado com sucesso.", patientProfile)

}

// @Summary Atualiza o perfil do Enfermeiro logado
// @Description Permite ao enfermeiro logado atualizar seu próprio perfil. Requer autenticação de Enfermeiro.
// @Description Campos protegidos (id, created_at, updated_at, password) não podem ser atualizados.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body object true "Campos para atualizar (JSON arbitrário). Ex: {\"bio\": \"Nova bio\", \"price\": 150.50}"
// @Success 200 {object} utils.SuccessNurseUpdateResponse "Usuário atualizado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou campo protegido"
// @Failure 401 {object} utils.ErrorResponse "Token inválido"
// @Failure 500 {object} utils.ErrorResponse "Erro interno ao atualizar"
// @Router /nurse/update [patch]
func (h *NurseHandler) UpdateNurseProfile(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	protectedFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"password":   true,
	}

	for key := range updates {
		if protectedFields[strings.ToLower(key)] {
			utils.SendErrorResponse(c, fmt.Sprintf("Campo(s) %s não pode ser atualizado.", key), http.StatusBadRequest)
			return
		}
	}

	user, err := h.nurseService.UpdateNurseFields(nurseId, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Usuário atualizado com sucesso.", user)
}

// @Summary Deleta o perfil do Enfermeiro logado
// @Description Deleta permanentemente a conta do enfermeiro logado. Requer senha atual para confirmação. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.DeleteAccountPasswordDto true "Senha atual para confirmar a exclusão"
// @Success 200 {object} utils.SuccessResponseNoData "Usuário deletado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou senha incorreta"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /nurse/delete [delete]
func (h *NurseHandler) DeleteNurseProfile(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	var deleteAccountPasswordDto dto.DeleteAccountPasswordDto
	if err := c.ShouldBindJSON(&deleteAccountPasswordDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.nurseService.DeleteNurse(nurseId, deleteAccountPasswordDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Usuário deletado com sucesso.", http.StatusOK)
}

// @Summary Informações de disponibilidade do Enfermeiro
// @Description Retorna as configurações de disponibilidade (horários, dias, bairros, etc.) do enfermeiro logado. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessAvailabilityResponse "Informações de disponibilidade carregadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar dados"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/availability [get]
func (h *NurseHandler) GetAvailabilityInfo(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseInfo, err := h.nurseService.GetAvailabilityInfo(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Informações de disponibilidade carregadas com sucesso.", nurseInfo)
}

// @Summary Informações detalhadas da Visita (Enfermeiro)
// @Description Retorna detalhes de uma visita específica e do paciente associado a ela. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita"
// @Success 200 {object} utils.SuccessNurseVisitInfoResponse "Informações de visita listadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao buscar visita"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/visit-info/{id} [get]
func (h *NurseHandler) GetNurseVisitInfo(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	visitInfo, err := h.nurseService.GetNurseVisitInfo(nurseId, visitId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Informações de visita para enfermeiro(a) listadas com sucesso.", visitInfo)

}

// @Summary Confirma a realização de um serviço (Visita)
// @Description Enfermeiro envia um código (fornecido pelo paciente) para confirmar que a visita foi concluída. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser confirmada"
// @Param payload body object true "Código de confirmação. Ex: {\"confirmation_code\": \"123456\"}"
// @Success 200 {object} utils.SuccessResponseNoData "Serviço completado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido, código obrigatório ou código incorreto"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/service-confirmation/{id} [patch]
func (h *NurseHandler) VisitServiceConfirmation(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	var requestBody map[string]interface{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Erro ao processar JSON"})
		return
	}

	confirmationCode, exists := requestBody["confirmation_code"]
	if !exists {
		c.JSON(400, gin.H{"error": "Campo confirmation_code é obrigatório"})
		return
	}

	codeStr, ok := confirmationCode.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "confirmation_code deve ser uma string"})
		return
	}

	err := h.nurseService.VisitServiceConfirmation(nurseId, visitId, codeStr)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Serviço completado com sucesso.", http.StatusOK)
}

// @Summary Fica offline ao fazer logout
// @Description Endpoint chamado pelo frontend no logout para garantir que o enfermeiro seja marcado como 'offline'. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessResponseNoData "Logout realizado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao atualizar status"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/offline [patch]
func (h *NurseHandler) TurnOfflineOnLogout(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	err := h.nurseService.TurnOfflineOnLogout(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Logout realizado com sucesso.", http.StatusOK)
}

// @Summary Rejeita uma visita pendente
// @Description Permite ao enfermeiro rejeitar uma visita que estava 'PENDING'. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser rejeitada"
// @Success 200 {object} utils.SuccessResponseNoData "Visita rejeitada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/reject-visit/{id} [patch]
func (h *NurseHandler) RejectVisit(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	err := h.nurseService.RejectVisit(nurseId, visitId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Visita rejeitada com sucesso.", http.StatusOK)
}

// @Summary Adiciona review para um Paciente
// @Description Permite ao enfermeiro avaliar um paciente (com nota e comentário) após uma visita. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita (para associar o review)"
// @Param payload body dto.ReviewDTO true "Dados de avaliação (Nota e Comentário)"
// @Success 200 {object} utils.SuccessResponseNoData "Review adicionada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou erro"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/review/{id} [post]
func (h *NurseHandler) AddReview(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	var reviewDto dto.ReviewDTO
	if err := c.ShouldBindJSON(&reviewDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.nurseService.AddReview(nurseId, visitId, reviewDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Review para paciente adicionada com sucesso.", http.StatusOK)
}

// @Summary Meu Perfil de Enfermeiro
// @Description Retorna o perfil completo do enfermeiro logado (idêntico ao /dashboard_info). Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessNurseProfileResponse "Perfil completo de enfermeiro(a) listado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar perfil"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/my-profile [get]
func (h *NurseHandler) GetMyNurseProfile(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseProfile, err := h.nurseService.GetMyNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

// @Summary Cria link de onboarding do Stripe
// @Description Gera um link único para o enfermeiro logado completar seu cadastro no Stripe. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessStripeOnboardingResponse "Link de onboarding criado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao criar link"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/stripe-onboarding [post]
func (h *NurseHandler) SetupStripeOnboarding(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	responseDto, err := h.nurseService.CreateStripeOnboardingLink(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Link de onboarding criado com sucesso.", responseDto)
}

// @Summary Adiciona prescrição a uma visita
// @Description Adiciona uma lista de prescrições a uma visita concluída. Requer autenticação de Enfermeiro.
// @Tags Nurse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita"
// @Param payload body dto.PrescriptionList true "Lista de prescrições"
// @Success 200 {object} utils.SuccessNurseVisitInfoResponse "Prescrição adicionada com sucesso (retorna a visita atualizada)"
// @Failure 400 {object} utils.ErrorResponse "ID inválido, JSON inválido ou erro"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Enfermeiro)"
// @Router /nurse/prescription/{id} [patch]
func (h *NurseHandler) AddPrescription(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	// permite que o campo 'reason' venha vazio
	var prescriptionList dto.PrescriptionList
	if err := c.ShouldBindJSON(&prescriptionList); err != nil {
		if err.Error() != "EOF" {
			utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
	}

	visit, err := h.nurseService.AddPrescription(nurseId, visitId, prescriptionList.PrescriptionList)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Prescrição adicionada com sucesso.", visit)
}
