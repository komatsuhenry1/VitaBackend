package chat

import (
	"log"
	"medassist/utils"

	"github.com/gin-gonic/gin"
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
