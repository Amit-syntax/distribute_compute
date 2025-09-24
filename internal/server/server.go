package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Action         string `json:"action"` // should be "register"
	ClientUsername string `json:"client_username"`

	// choices{system_info_update, code_run}
	MessageType string `json:"message_type"`

	Content map[string]any `json:"content"`
}

type RegisterMessage struct {
	Action         string `json:"action"` // should be "register"
	ClientUsername string `json:"client_username"`
	JoineeType     string `json:"joinee_type"` // choices{worker,consumer}
}

type Client struct {
	IP         string `json:"ip"`
	Username   string `json:"username"`
	JoineeType string `json:"joinee_type"`
	conn       *websocket.Conn
	hub        *Hub
}

type Hub struct {
	clients map[*Client]bool
	mu      *sync.RWMutex
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client.hub = h
	h.clients[client] = true
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client.hub = h
	h.clients[client] = false
}

var hub = &Hub{
	clients: make(map[*Client]bool),
	mu:      &sync.RWMutex{},
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (client *Client) readBulk() {

	defer func() {
		client.hub.Unregister(client)
		client.conn.Close()
	}()

	for {
		_, msgByte, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{}
		if err = json.Unmarshal(msgByte, msg); err != nil {
			log.Printf("error unmarshaling: %v", err)
			continue
		}
		// TODO: Process message
		log.Printf("message: %v", msg)
	}

}

func HandleClientConnHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Wait for client registration message
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	_, msgByte, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
	}

	msg := &RegisterMessage{}
	if err = json.Unmarshal(msgByte, msg); err != nil {
		log.Printf("error unmarshaling: %v", err)
	}

	client := &Client{
		Username:   msg.ClientUsername,
		JoineeType: msg.JoineeType,
	}

	hub.Register(client)

	log.Printf("clients: %d", len(hub.clients))

}

func RunJobHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: recieve job connection request (only for already connected consumers)

}
