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

// @Summary Lista todos os enfermeiros (para agendar)
// @Description Retorna todos os enfermeiros na cidade do paciente logado. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessAllNursesResponse "Enfermeiros listados com sucesso"
// @Failure 401 {object} utils.ErrorResponse "Token inválido"
// @Failure 500 {object} utils.ErrorResponse "Erro ao buscar enfermeiros"
// @Router /user/all_nurses [get]
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

// @Summary Exibe um arquivo (ex: imagem de perfil)
// @Description Retorna um arquivo do GridFS (como uma imagem de perfil) para ser exibido 'inline' no navegador. (Endpoint atualmente público).
// @Tags User
// @Produce image/png
// @Produce image/jpeg
// @Produce application/octet-stream
// @Param id path string true "ID do Arquivo (GridFS ObjectID)"
// @Success 200 {file} file "A imagem ou arquivo"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou Arquivo não encontrado"
// @Router /user/file/{id} [get]
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

// @Summary Envia mensagem de contato
// @Description Endpoint público para enviar uma mensagem (dúvida, sugestão, reclamação) para a central de contato.
// @Tags User
// @Accept json
// @Produce json
// @Param payload body dto.ContactUsDTO true "Dados da mensagem de contato"
// @Success 200 {object} utils.SuccessResponseNoData "Mensagem enviada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida (campos obrigatórios faltando)"
// @Failure 500 {object} utils.ErrorResponse "Erro interno ao enviar mensagem"
// @Router /user/contact [post]
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

// @Summary Perfil público do Enfermeiro
// @Description Retorna o perfil detalhado de um enfermeiro específico (para agendamento). Requer autenticação de Paciente ou Enfermeiro.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Enfermeiro"
// @Success 200 {object} utils.SuccessNurseProfileResponse "Perfil completo de enfermeiro(a) listado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao buscar perfil"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/nurse/{id} [get]
func (h *UserHandler) GetNurseProfile(c *gin.Context) {
	nurseId := c.Param("id")

	nurseProfile, err := h.userService.GetNurseProfile(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Perfil completo de enfermeiro(a) listado com sucesso", nurseProfile)
}

// @Summary Solicita uma visita agendada
// @Description Cria uma nova solicitação de visita agendada para um enfermeiro específico. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.CreateVisitDto true "Detalhes da visita agendada"
// @Success 200 {object} utils.SuccessResponseNoData "Visita agendada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou erro na solicitação"
// @Failure 401 {object} utils.ErrorResponse "Token inválido"
// @Router /user/visit [post]
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

// @Summary Lista todas as visitas do Paciente
// @Description Retorna um histórico de todas as visitas (pendentes, concluídas, etc.) do paciente logado. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessVisitsListResponse "Lista de visitas listadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao buscar visitas"
// @Failure 401 {object} utils.ErrorResponse "Token inválido"
// @Router /user/visits [get]
func (h *UserHandler) GetAllVisits(c *gin.Context) {
	patientId := utils.GetUserId(c)

	visits, err := h.userService.FindAllVisits(patientId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Lista de visitas listadas com sucesso.", visits)

}

// @Summary Atualiza o perfil do Paciente
// @Description Permite ao paciente logado atualizar seu próprio perfil. Requer autenticação de Paciente.
// @Description Campos protegidos (id, created_at, updated_at) não podem ser atualizados.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body object true "Campos para atualizar (JSON arbitrário). Ex: {\"name\": \"Novo Nome\", \"phone\": \"11999998888\"}"
// @Success 200 {object} utils.SuccessUserTypeResponse "Usuário atualizado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou campo protegido"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/update [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	patientId := utils.GetUserId(c)

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

// @Summary Deleta o perfil do Paciente
// @Description Deleta permanentemente a conta do paciente logado. Requer senha atual para confirmação. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.DeleteAccountPasswordDto true "Senha atual para confirmação"
// @Success 200 {object} utils.SuccessResponseNoData "Usuário deletado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou senha incorreta"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/delete [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	patientId := utils.GetUserId(c)

	var deleteAccountPasswordDto dto.DeleteAccountPasswordDto
	if err := c.ShouldBindJSON(&deleteAccountPasswordDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.DeleteUser(patientId, deleteAccountPasswordDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário deletado com sucesso.", http.StatusOK)
}

// @Summary Confirma a conclusão de um serviço (Paciente)
// @Description O paciente confirma que o serviço da visita foi concluído. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser concluída"
// @Success 200 {object} utils.SuccessResponseNoData "Serviço concluído com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro na confirmação"
// @Failure 401 {object} utils.ErrorResponse "Token inválido"
// @Router /user/visit/{id} [patch]
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

// @Summary Lista enfermeiros online
// @Description Retorna todos os enfermeiros que estão online e na cidade do paciente logado. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessAllNursesResponse "Lista de enfermeiros online listada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao buscar enfermeiros"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Paciente)"
// @Router /user/online_nurses [get]
func (h *UserHandler) GetOnlineNurses(c *gin.Context) {
	userId := utils.GetUserId(c)

	response, err := h.userService.GetOnlineNurses(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Lista de enfermeiros online listada com sucesso.", response)
}

// @Summary Informações detalhadas da Visita (Paciente)
// @Description Retorna detalhes de uma visita específica e do enfermeiro associado a ela. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita"
// @Success 200 {object} utils.SuccessPatientVisitInfoResponse "Informações de visita listadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/visit-info/{id} [get]
func (h *UserHandler) GetPatientVisitInfo(c *gin.Context) {
	patientId := utils.GetUserId(c)

	visitId := c.Param("id")

	visitInfo, err := h.userService.GetPatientVisitInfo(patientId, visitId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Informações de visita para paciente listadas com sucesso.", visitInfo)
}

// @Summary Adiciona review para um Enfermeiro
// @Description Permite ao paciente avaliar um enfermeiro (com nota e comentário) após uma visita. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita (para associar o review)"
// @Param payload body dto.ReviewDTO true "Dados da avaliação (Nota e Comentário)"
// @Success 200 {object} utils.SuccessResponseNoData "Review adicionada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido ou erro"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/review/{id} [post]
func (h *UserHandler) AddReview(c *gin.Context) {
	userId := utils.GetUserId(c)

	visitId := c.Param("id")

	var reviewDto dto.ReviewDTO
	if err := c.ShouldBindJSON(&reviewDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.AddReview(userId, visitId, reviewDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Review para enfermeiro(a) adicionada com sucesso.", http.StatusOK)
}

// @Summary Solicita uma visita imediata
// @Description Cria uma nova solicitação de visita imediata para um enfermeiro (que deve estar online). Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.ImmediateVisitDTO true "Detalhes da visita imediata"
// @Success 200 {object} utils.SuccessImmediateVisitResponse "Visita imediata solicitada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido, enfermeiro offline ou erro na solicitação"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Paciente)"
// @Router /user/immediate-visit [post]
func (h *UserHandler) ImmediateVisitSolicitation(c *gin.Context) {
	patientId := utils.GetUserId(c)

	var immediateVisitDto dto.ImmediateVisitDTO
	if err := c.ShouldBindJSON(&immediateVisitDto); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.userService.ImmediateVisitSolicitation(patientId, immediateVisitDto)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Visita imediata solicitada com sucesso.", gin.H{"patient_id": userId})
}

// @Summary Meu Perfil de Paciente
// @Description Retorna o perfil completo do paciente logado, incluindo dados privados. Requer autenticação de Paciente.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessPatientProfileResponsee "Perfil do paciente listado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar perfil"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido"
// @Router /user/my-profile [get]
func (h *UserHandler) GetMyUserProfile(c *gin.Context){

	userId := utils.GetUserId(c)

	patientProfile, err := h.userService.GetPatientProfile(userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Perfil do paciente listado com sucesso.", patientProfile)

}
