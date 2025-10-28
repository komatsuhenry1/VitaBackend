package chat

import (
	"log" // Adicionado para logs
	"medassist/internal/repository"
)

type Hub struct {
	// 1. ALTERADO: Mapeia UserID (string) para *Client
	// Isso permite encontrar um cliente específico pelo ID
	clients map[string]*Client

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	msgRepo    repository.MessageRepository
}

func NewHub(msgRepo repository.MessageRepository) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// Inicializa o mapa modificado
		clients: make(map[string]*Client),
		msgRepo: msgRepo,
	}
}

// 2. NOVO MÉTODO: SendToNurse
// Envia uma mensagem para um cliente específico baseado no UserID.
// Retorna true se o cliente foi encontrado e a mensagem enfileirada, false caso contrário.
func (h *Hub) SendToNurse(userID string, message []byte) bool {
	// Procura o cliente no mapa usando o UserID
	client, ok := h.clients[userID]
	if !ok {
		// Cliente não está conectado ao WebSocket
		log.Printf("[Hub SendToNurse] Tentativa de enviar para usuário offline ou não conectado: %s", userID)
		return false
	}

	// Tenta enviar a mensagem para o canal 'send' do cliente
	select {
	case client.send <- message:
		log.Printf("[Hub SendToNurse] Mensagem enviada com sucesso para: %s", userID)
		return true // Mensagem enfileirada
	default:
		// O canal 'send' do cliente está cheio ou fechado (cliente lento/desconectado)
		// Remove o cliente do hub para evitar tentar enviar novamente
		log.Printf("[Hub SendToNurse] Canal do usuário %s cheio ou fechado. Removendo cliente.", userID)
		delete(h.clients, userID) // Usa o ID para remover
		close(client.send)
		// IMPORTANTE: Aqui NÃO chamamos SetNurseOffline, pois este Hub é genérico.
		// A lógica de SetNurseOffline deve ficar no hub específico de visitas (se você o criar)
		// ou ser tratada na desconexão (unregister).
		return false // Falha ao enviar
	}
}

func (h *Hub) Run() {
	for {
		select {
		// Caso um novo cliente se conecte
		case client := <-h.register:
			// 3. ALTERADO: Usa UserID como chave
			log.Printf("[Hub Run] Registrando cliente: ID=%s, Nome=%s, Role=%s", client.UserID, client.Name, client.Role)
			h.clients[client.UserID] = client
			// NOTA: A lógica de SetNurseOnline NÃO entra aqui, pois este Hub é para CHAT.
			// Se você quiser que a conexão ao chat marque o enfermeiro como online,
			// você precisaria injetar o NurseService aqui e chamar SetNurseOnline.
			// Mas isso mistura as responsabilidades do chat e do status de visita.

		// Caso um cliente se desconecte
		case client := <-h.unregister:
			// 4. ALTERADO: Verifica e remove usando UserID
			if _, ok := h.clients[client.UserID]; ok {
				log.Printf("[Hub Run] Desregistrando cliente: ID=%s", client.UserID)
				delete(h.clients, client.UserID)
				close(client.send)
				// NOTA: A lógica de SetNurseOffline NÃO entra aqui pelo mesmo motivo acima.
				// A desconexão do CHAT não deveria necessariamente marcar o enfermeiro como
				// indisponível para VISITAS, a menos que seja essa sua regra de negócio.
			}

		// Caso chegue uma nova mensagem para ser enviada a TODOS (Broadcast - Lógica de Chat)
		case message := <-h.broadcast:
			// Esta parte permanece igual, envia para todos no mapa
			log.Printf("[Hub Run] Broadcasting mensagem para %d clientes", len(h.clients))
			for userID, client := range h.clients { // Itera sobre o mapa modificado
				select {
				case client.send <- message:
					// Enviado com sucesso
				default:
					// Cliente lento/desconectado durante o broadcast
					log.Printf("[Hub Run - Broadcast] Canal do usuário %s cheio ou fechado. Removendo cliente.", userID)
					delete(h.clients, userID)
					close(client.send)
				}
			}
		}
	}
}
