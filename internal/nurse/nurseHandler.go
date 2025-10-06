package nurse

import (
	"medassist/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"medassist/internal/nurse/dto"
	"fmt"
)

type NurseHandler struct {
	nurseService NurseService
}

func NewNurseHandler(nurseService NurseService) *NurseHandler {
	return &NurseHandler{nurseService: nurseService}
}

func (h *NurseHandler) NurseDashboard(c *gin.Context){
	nurseId := utils.GetUserId(c)

	dashboardData, err := h.nurseService.NurseDashboardData(nurseId)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Dados de enfermeiro carregados com sucesso.", dashboardData)
}

func (h *NurseHandler) ChangeOnlineNurse(c *gin.Context){
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

func (h *NurseHandler) GetAllVisits(c *gin.Context){
	//pega o id pelo token
	nurseId := utils.GetUserId(c)

	visits, err := h.nurseService.GetAllVisits(nurseId)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Visitas [PENDENTE/MARCADAS/CONCLUIDAS] listadas com sucesso.", visits)
}

func (h *NurseHandler) ConfirmOrCancelVisit(c *gin.Context){
	nurseId := utils.GetUserId(c) // pega o id do user pela req

	fmt.Println("nurseiD", nurseId)

	visitId := c.Param("id")
	fmt.Println("antes binfing")
	var reason dto.CancelReason
	if err := c.ShouldBindJSON(&reason); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("antes service")
	response, err := h.nurseService.ConfirmOrCancelVisit(nurseId, visitId, reason.Reason)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}	
	fmt.Println("response", response)
	
	utils.SendSuccessResponse(c, "Status da visita alterado com sucesso.", response)
}

func (h *NurseHandler) GetPatientProfile(c *gin.Context){

	patientId := c.Param("id")

	patientProfile, err := h.nurseService.GetPatientProfile(patientId)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)	
	}

	utils.SendSuccessResponse(c, "Perfil do paciente listado com sucesso.", patientProfile)

}
