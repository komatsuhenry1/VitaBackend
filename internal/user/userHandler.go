package user

import (
	"medassist/utils"
	"net/http"
	"medassist/internal/user/dto"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UserDashboard(c *gin.Context){
	utils.SendSuccessResponse(c, "user dashboard", http.StatusOK)
}

func (h *UserHandler) GetAllNurses(c *gin.Context){
	nurses, err := h.userService.GetAllNurses()
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
        utils.SendErrorResponse(c, "Arquivo não encontrado", http.StatusNotFound)
        return
    }

    c.Header("Content-Disposition", "inline; filename=\""+fileData.Filename+"\"")
    c.Data(http.StatusOK, fileData.ContentType, fileData.Data)
}

func (h *UserHandler) ContactUsMessage(c *gin.Context){

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

func (h *UserHandler) GetNurseProfile(c *gin.Context){
	nurseId := c.Param("id")
	
	nurseProfile, err := h.userService.GetNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

func (h *UserHandler) CreateVisit(c *gin.Context){
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
	
	err := h.userService.CreateVisit(patientId, createVisitDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Visita agendada com sucesso.", http.StatusOK)

}