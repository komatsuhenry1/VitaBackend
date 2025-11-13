package chat

import (
	"medassist/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medassist/internal/model"
)

type ChatHandler struct {
	msgRepo repository.MessageRepository
}

func NewChatHandler(msgRepo repository.MessageRepository) *ChatHandler {
	return &ChatHandler{
		msgRepo: msgRepo,
	}
}

// @Summary Histórico de mensagens com um enfermeiro
// @Description Retorna todas as mensagens entre o usuário logado (seja Paciente ou Enfermeiro) e um enfermeiro específico.
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param nurseId path string true "ID do Enfermeiro (com quem a conversa acontece)"
// @Success 200 {object} utils.SuccessMessagesResponse "Histórico de mensagens retornado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID do enfermeiro inválido"
// @Failure 401 {object} utils.ErrorResponse "Usuário não autenticado"
// @Failure 500 {object} utils.ErrorResponse "Erro ao buscar mensagens"
// @Router /chat/messages/{nurseId} [get]
func (h *ChatHandler) GetMessagesHistory(c *gin.Context) {
	nurseIdStr := c.Param("nurseId")
	nurseId, err := primitive.ObjectIDFromHex(nurseIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID do enfermeiro inválido"})
		return
	}

	userIDCtx, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Usuário não autenticado"})
		return
	}
	userIDStr := userIDCtx.(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "ID de usuário inválido"})
		return
	}

	messages, err := h.msgRepo.FindMessagesBetween(userID, nurseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Erro ao buscar mensagens"})
		return
	}

	if messages == nil {
		messages = make([]model.Message, 0)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
	})
}