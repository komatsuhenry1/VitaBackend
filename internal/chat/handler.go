package chat

import (
	"log"
	"medassist/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medassist/internal/chat/dto"
	"net/http"
)

func ServeWs(hub *Hub, c *gin.Context) {
	// 1. Pega o token da URL
	tokenStr := c.Query("token")
	if tokenStr == "" {
		log.Println("Erro no WebSocket: Token não fornecido")
		return // Fecha a conexão silenciosamente
	}

	// 2. Valida o token para obter os dados do usuário
	claims, err := utils.ValidateToken(tokenStr) // Use sua função de validação de token
	if err != nil {
		log.Printf("Erro no WebSocket: Token inválido - %v", err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		// Adiciona os dados do usuário ao cliente
		// CORREÇÃO: Acessando os dados como um mapa
		UserID: claims["sub"].(string), // A chave padrão para ID de usuário no JWT é "sub"
		Name:   claims["name"].(string),
		Role:   claims["role"].(string),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *ChatHandler) GetNurseConversations(c *gin.Context) {
	// Pega o ID do enfermeiro logado a partir do contexto (do middleware)
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

	// Chama a nova função do repositório
	conversations, err := h.msgRepo.GetConversationsForNurse(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Erro ao buscar conversas"})
		return
	}

	if conversations == nil {
		conversations = make([]dto.ConversationDTO, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    conversations,
	})
}

func (h *ChatHandler) GetPatientConversations(c *gin.Context) {
    // Pega o ID do paciente logado a partir do contexto
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

    // Chama a nova função do repositório para pacientes
    conversations, err := h.msgRepo.GetConversationsForPatient(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Erro ao buscar conversas"})
        return
    }

    if conversations == nil {
        conversations = make([]dto.ConversationDTO, 0)
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    conversations,
    })
}
