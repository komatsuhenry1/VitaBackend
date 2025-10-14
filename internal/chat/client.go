package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client agora tem os dados do usuário
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	// Novos campos para identificar o usuário
	UserID string
	Name   string
	Role   string
}

// upgrader é usado para promover uma conexão HTTP para WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Mensagem que será enviada via WebSocket (com todos os campos)
type WebSocketMessage struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	SenderRole string `json:"sender_role"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

// readPump agora enriquece a mensagem antes de enviá-la
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// ...
			break
		}

		// Cria a mensagem completa com os dados do remetente (do próprio client)
		fullMessage := WebSocketMessage{
			ID:         time.Now().String(), // ID temporário, o ideal seria do DB
			SenderID:   c.UserID,
			SenderName: c.Name,
			SenderRole: c.Role,
			Message:    string(message), // O conteúdo da mensagem que chegou
			Timestamp:  time.Now().Format(time.RFC3339),
		}

		// Converte a mensagem completa para JSON
		jsonMessage, err := json.Marshal(fullMessage)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			continue
		}

		// Envia a mensagem completa (enriquecida) para o hub
		c.hub.broadcast <- jsonMessage
	}
}

// writePump envia mensagens do hub para o websocket do cliente
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// O hub fechou o canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
