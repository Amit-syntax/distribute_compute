package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientRegistrationMsg struct {
  Name string `json:"name"` 
}

type Client struct {
	Name string `json:"name"`
}

type Hub struct {
	clients map[*Client]bool
	mu *sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

var hub = Hub{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (client *Client) readBulk() {
	
}

func handleWebsocketConn(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Wait for registration message
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	_, messageBytes, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read registration message: %v", err)
		conn.Close()
		return
	}

	var regMsg ClientRegistrationMsg
	if err := json.Unmarshal(messageBytes, &regMsg); err != nil {
		log.Printf("Failed to unmarshal registration message: %v", err)
		conn.Close()
		return
	}
}


func NewServer() {
	// TODO:
	// run a websocket server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocketConn(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print(w, "WebSocket Server is running!\nConnect to ws://localhost:8080/ws\nCheck status at http://localhost:8080/status")
	})

	log.Println("WebSocket server starting on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("Status endpoint: http://localhost:8080/status")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

