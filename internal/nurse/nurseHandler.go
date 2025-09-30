package nurse

import (
	"medassist/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type NurseHandler struct {
	nurseService NurseService
}

func NewNurseHandler(nurseService NurseService) *NurseHandler {
	return &NurseHandler{nurseService: nurseService}
}

func (h *NurseHandler) NurseDashboard(c *gin.Context){
	utils.SendSuccessResponse(c, "nurse dashboard", http.StatusOK)
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

func (h *NurseHandler) ConfirmVisit(c *gin.Context){
	utils.SendSuccessResponse(c, "confirm visit", http.StatusOK)
}
