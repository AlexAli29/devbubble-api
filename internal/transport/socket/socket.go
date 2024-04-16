package socket

import (
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/validator"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type SocketHandler struct {
	validate      *validator.Validator
	updater       websocket.Upgrader
	jwtMiddleware func(next http.Handler) http.Handler
	log           *slog.Logger
	sync.RWMutex
	handlers map[string]EventHandler
	clients  ClientList
	queries  *db.Queries
}

func NewSocketHandler(validate *validator.Validator, jwtMiddleware func(next http.Handler) http.Handler, updater websocket.Upgrader, log *slog.Logger, queries *db.Queries) *SocketHandler {

	h := &SocketHandler{
		validate:      validate,
		jwtMiddleware: jwtMiddleware,
		updater:       updater,
		log:           log,
		clients:       make(ClientList),
		handlers:      make(map[string]EventHandler),
		queries:       queries,
	}
	h.setupEventHandlers()
	return h
}

func (h *SocketHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.With(h.jwtMiddleware).Get("/", h.handleWebSocket)

	return r
}

func (h *SocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("UserId").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	conn, err := h.updater.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(err.Error())
		return
	}
	client := NewClient(conn, h)
	h.addClient(userId, client)
	go client.readMessages()
	go client.writeMessages()
}

func (h *SocketHandler) SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var message core.Message
	if err := json.Unmarshal(event.Payload, &message); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data

	outgoingEvent.Type = EventNewPrivateMessage
	// Broadcast to all other Clients

	for _, client := range c.handler.clients {
		// Only send to clients inside the same chatroom
		client.egress <- outgoingEvent
		// if client.room == c.room {

		// }

	}
	return nil
}

func (h *SocketHandler) setupEventHandlers() {
	h.handlers[EventNewPrivateMessage] = h.SendMessageHandler
	h.handlers[EventChangeRoom] = ChatRoomHandler
}

func (h *SocketHandler) routeEvent(event Event, c *Client) error {

	// Check if Handler is present in Map
	if handler, ok := h.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// addClient will add clients to our clientList
func (m *SocketHandler) addClient(id string, client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[id] = client
}

// removeClient will remove the client and clean up
func (m *SocketHandler) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Поиск ключа для клиента в карте
	for key, value := range m.clients {
		if value == client {
			// Закрытие соединения
			client.conn.Close()
			// Удаление клиента из карты
			delete(m.clients, key)
			break
		}
	}
}
