package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Amit-syntax/distribute_compute/internal/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id         string            `json:"id"`
	Username   string            `json:"username"`
	JoineeType common.JoineeType `json:"joinee_type"`
	conn       *websocket.Conn
	hub        *ClientHub
}

type ClientHub struct {
	clients map[*Client]bool
	mu      *sync.RWMutex
}

func (h *ClientHub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client.hub = h
	h.clients[client] = true
}

func (h *ClientHub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client.hub = h
	h.clients[client] = false
}

func (h *ClientHub) GetClientById(uuid string) *Client {

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.Id == uuid && h.clients[client] {
			return client
		}
	}
	return nil
}

var hub = &ClientHub{
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
		
		msg := &common.Message{}
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

	msg := &common.RegisterMessage{}
	if err = json.Unmarshal(msgByte, msg); err != nil {
		log.Printf("error unmarshaling: %v", err)
	}

	client := &Client{
		Id:         uuid.New().String(),
		Username:   msg.ClientUsername,
		JoineeType: msg.JoineeType,
		conn:       conn,
	}

	hub.Register(client)
	go client.readBulk()

	log.Printf("clients: %d", len(hub.clients))

}
