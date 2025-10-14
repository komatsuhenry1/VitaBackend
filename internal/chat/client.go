package chat

import (
	"encoding/json"
	"log"
	"medassist/internal/model"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
type ClientMessage struct {
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
}

// Substitua sua função readPump inteira por esta versão corrigida
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			// ...
			break
		}

		// 1. Decodifica a mensagem que chegou do cliente
		var clientMsg ClientMessage
		if err := json.Unmarshal(rawMessage, &clientMsg); err != nil {
			log.Printf("error unmarshaling client message: %v", err)
			continue
		}

		// 2. Converte os IDs de string para ObjectID
		senderID, _ := primitive.ObjectIDFromHex(c.UserID)
		receiverID, _ := primitive.ObjectIDFromHex(clientMsg.ReceiverID)

		// 3. Cria uma instância do modelo de MENSAGEM PARA O BANCO DE DADOS
		dbMessage := &model.Message{
			SenderID:   senderID,
			ReceiverID: receiverID,
			Content:    clientMsg.Message,
			Read:       false,
		}

		// 4. SALVA A MENSAGEM NO BANCO DE DADOS
		err = c.hub.msgRepo.Save(dbMessage)
		if err != nil {
			log.Printf("error saving message to db: %v", err)
			continue
		}
		// Após salvar, dbMessage agora contém o ID e o Timestamp gerados pelo banco

		// 5. Cria a mensagem para transmitir via WebSocket com os dados salvos
		broadcastMessage := WebSocketMessage{
			ID:         dbMessage.ID.Hex(),
			SenderID:   dbMessage.SenderID.Hex(),
			SenderName: c.Name,
			SenderRole: c.Role,
			Message:    dbMessage.Content,
			Timestamp:  dbMessage.Timestamp.Format(time.RFC3339),
		}

		jsonMessage, _ := json.Marshal(broadcastMessage)
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
