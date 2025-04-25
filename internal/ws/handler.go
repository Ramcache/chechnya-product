package ws

import (
	"chechnya-product/internal/middleware"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections
// @Summary Подключение к WebSocket для уведомлений о заказах
// @Description Устанавливает WebSocket-соединение. Админы получают уведомления о новых заказах.
// @Tags WebSocket
// @Produce json
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/orders [get]
func (h *Hub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Пытаемся извлечь userID и роль
	userID := middleware.GetUserIDOrZero(r)
	role := "guest"
	if userID > 0 {
		role = middleware.GetUserRole(r)
	}

	client := &Client{
		ID:   userID,
		Role: role,
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  h,
	}

	h.register <- client

	go client.readPump()
	go client.writePump()
}

// HandleAnnouncementConnections
// @Summary WebSocket подключение для объявлений
// @Tags WebSocket
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/announcements [get]
func (h *Hub) HandleAnnouncementConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	userID := middleware.GetUserIDOrZero(r)
	role := "guest"
	if userID > 0 {
		role = middleware.GetUserRole(r)
	}
	client := &Client{
		ID:   userID,
		Role: role,
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  h,
	}
	h.register <- client
	go client.readPump()
	go client.writePump()
}

// readPump — читает сообщения от клиента (игнорируем)
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
	}()
	for {
		if _, _, err := c.Conn.NextReader(); err != nil {
			break
		}
	}
}

// writePump — отправляет сообщения клиенту
func (c *Client) writePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
