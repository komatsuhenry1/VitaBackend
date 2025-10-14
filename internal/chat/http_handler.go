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