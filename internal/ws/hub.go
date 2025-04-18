package ws

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"sync"
)

// Client — подключённый клиент WebSocket
type Client struct {
	ID   int           // userID или 0 для гостей
	Role string        // "admin", "user" или "guest"
	Conn WebSocketConn // интерфейс вместо *websocket.Conn для тестируемости
	Send chan []byte   // канал для отправки сообщений клиенту
	Hub  *Hub          // ссылка на хаб
}

// WebSocketConn — интерфейс для WebSocket-соединения
type WebSocketConn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
	NextReader() (messageType int, r io.Reader, err error)
}

// OrderMessage — формат сообщения для рассылки
type OrderMessage struct {
	Type  string       `json:"type"`  // тип сообщения, например: "new_order"
	Order models.Order `json:"order"` // заказ
}

// Hub — управляет всеми WebSocket-клиентами
type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcastCh chan OrderMessage
	mu          sync.Mutex
	logger      *zap.Logger // 👈 логгер
}

// NewHub создаёт новый экземпляр Hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcastCh: make(chan OrderMessage),
		logger:      logger,
	}
}

// Run запускает основной цикл хаба
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Info("WebSocket client connected",
				zap.Int("user_id", client.ID),
				zap.String("role", client.Role),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				client.Conn.Close()
				h.logger.Info("WebSocket client disconnected",
					zap.Int("user_id", client.ID),
					zap.String("role", client.Role),
				)
			}
			h.mu.Unlock()

		case msg := <-h.broadcastCh:
			data, _ := json.Marshal(msg)

			h.mu.Lock()
			sent := 0
			for client := range h.clients {
				select {
				case client.Send <- data:
					sent++
				default:
					close(client.Send)
					delete(h.clients, client)
					client.Conn.Close()
				}
			}
			h.mu.Unlock()

			h.logger.Info("WebSocket broadcast",
				zap.String("type", msg.Type),
				zap.Int("order_id", msg.Order.ID),
				zap.Int("clients_sent", sent),
			)
		}
	}
}

// BroadcastNewOrder — рассылает заказ пользователю и всем админам
func (h *Hub) BroadcastNewOrder(order models.Order) {
	h.broadcastCh <- OrderMessage{
		Type:  "new_order",
		Order: order,
	}
}
