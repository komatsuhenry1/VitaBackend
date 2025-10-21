package user

import (
	"medassist/internal/user/dto"
	"medassist/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UserDashboard(c *gin.Context) {
	utils.SendSuccessResponse(c, "user dashboard", http.StatusOK)
}

func (h *UserHandler) GetAllNurses(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	nurses, err := h.userService.GetAllNurses(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Enfermeiros listados com sucesso.", nurses)
}

func (h *UserHandler) GetFileByID(c *gin.Context) {
	fileIDStr := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(fileIDStr)
	if err != nil {
		utils.SendErrorResponse(c, "ID de arquivo inválido", http.StatusBadRequest)
		return
	}

	fileData, err := h.userService.GetFileByID(c.Request.Context(), objectID)
	if err != nil {
		utils.SendErrorResponse(c, "Arquivo não encontrado", http.StatusBadRequest)
		return
	}

	c.Header("Content-Disposition", "inline; filename=\""+fileData.Filename+"\"")
	c.Data(http.StatusOK, fileData.ContentType, fileData.Data)
}

func (h *UserHandler) ContactUsMessage(c *gin.Context) {

	var contactUsDto dto.ContactUsDTO
	if err := c.ShouldBindJSON(&contactUsDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.ContactUsMessage(contactUsDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
	}

	utils.SendSuccessResponse(c, "Mensagem de contato para central enviada com sucesso.", http.StatusOK)
}

func (h *UserHandler) GetNurseProfile(c *gin.Context) {
	nurseId := c.Param("id")

	nurseProfile, err := h.userService.GetNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

func (h *UserHandler) VisitSolicitation(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	var createVisitDto dto.CreateVisitDto
	if err := c.ShouldBindJSON(&createVisitDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.VisitSolicitation(patientId, createVisitDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Visita agendada com sucesso.", http.StatusOK)
}

func (h *UserHandler) GetAllVisits(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	visits, err := h.userService.FindAllVisits(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Lista de visitas listadas com sucesso.", visits)

}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
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
	}

	for key := range updates {
		if protectedFields[strings.ToLower(key)] {
			utils.SendErrorResponse(c, fmt.Sprintf("Campo(s) %s não pode ser atualizado.", key), http.StatusBadRequest)
			return
		}
	}

	user, err := h.userService.UpdateUser(patientId, updates)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Usuário atualizado com sucesso.", user)

}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	err := h.userService.DeleteUser(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Usuário deletado com sucesso.", http.StatusOK)
}

func (h *UserHandler) ConfirmVisitService(c *gin.Context) {
	visitId := c.Param("id")

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	patientId, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	err := h.userService.ConfirmVisitService(visitId, patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Serviço concluído com sucesso.", http.StatusOK)
}
