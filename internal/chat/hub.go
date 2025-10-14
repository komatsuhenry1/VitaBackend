// Será o "gerente" do nosso chat. Ele vai controlar todos os clientes conectados e transmitir as mensagens para todos.

package chat

import "medassist/internal/repository" // <-- ADICIONE ESTA IMPORTAÇÃO

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	msgRepo    repository.MessageRepository // <-- ADICIONE ESTA LINHA
}

// NewHub cria uma nova instância do Hub.
func NewHub(msgRepo repository.MessageRepository) *Hub { // <-- ALTERE O PARÂMETRO
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		msgRepo:    msgRepo, // <-- ATRIBUA AQUI
	}
}

// Run inicia o processamento de ações do Hub.
// Esta função deve ser executada em uma goroutine separada.
func (h *Hub) Run() {
	for {
		select {
		// Caso um novo cliente se conecte
		case client := <-h.register:
			h.clients[client] = true

		// Caso um cliente se desconecte
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		// Caso chegue uma nova mensagem para ser enviada a todos
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				// Se o buffer do canal de envio estiver cheio, assumimos que o cliente está lento
				// ou desconectado, então não enviamos a mensagem.
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
