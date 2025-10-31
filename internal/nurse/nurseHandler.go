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

func (h *NurseHandler) NurseDashboard(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	dashboardData, err := h.nurseService.NurseDashboardData(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Dados de enfermeiro carregados com sucesso.", dashboardData)
}

func (h *NurseHandler) NurseDashboardData(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseProfile, err := h.nurseService.GetNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

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

func (h *NurseHandler) GetPatientProfile(c *gin.Context) {

	patientId := c.Param("id")

	patientProfile, err := h.nurseService.GetPatientProfile(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Perfil do paciente listado com sucesso.", patientProfile)

}

func (h *NurseHandler) UpdateNurseProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	nurseId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

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

func (h *NurseHandler) GetAvailabilityInfo(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseInfo, err := h.nurseService.GetAvailabilityInfo(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Informações de disponibilidade carregadas com sucesso.", nurseInfo)
}

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

func (h *NurseHandler) TurnOfflineOnLogout(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	err := h.nurseService.TurnOfflineOnLogout(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Logout realizado com sucesso.", http.StatusOK)
}

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

func (h *NurseHandler) AddReview(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	visitId := c.Param("id")

	fmt.Println("===")

	var reviewDto dto.ReviewDTO
	if err := c.ShouldBindJSON(&reviewDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("===")
	err := h.nurseService.AddReview(nurseId, visitId, reviewDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Review para paciente adicionada com sucesso.", http.StatusOK)
}

func (h *NurseHandler) GetMyNurseProfile(c *gin.Context) {
	nurseId := utils.GetUserId(c)

	nurseProfile, err := h.nurseService.GetMyNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}
