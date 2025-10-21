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
	nurseId := c.Param("id")

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
	fmt.Println("antes binfing")
	var reason dto.CancelReason
	if err := c.ShouldBindJSON(&reason); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
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

	err := h.nurseService.DeleteNurse(nurseId)
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
