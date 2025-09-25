package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionMsg struct {
	PipModules []string `json:"pip_modules"`
	ClientId   string   `json:"client_id"`
}

type SessionConn struct {
	SessionId  string
	Client     *Client
	PipModules []string
	CreatedAt  time.Time
}

type SessionHub struct {
	sessions map[string]*SessionConn
	mu       *sync.RWMutex
}

func NewSessionHub() *SessionHub {
	return &SessionHub{
		sessions: make(map[string]*SessionConn),
		mu:       &sync.RWMutex{},
	}
}

func (h *SessionHub) Register(session *SessionConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sessions[session.SessionId] = session
}

var sessionHub = NewSessionHub()

func SessionHandler(w http.ResponseWriter, r *http.Request) {

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
		return
	}

	msg := &SessionMsg{}
	if err = json.Unmarshal(msgByte, msg); err != nil {
		log.Printf("error unmarshaling: %v", err)
		return
	}

	client := hub.GetClientById(msg.ClientId)
	if client == nil {
		log.Printf("Client with ID %s not found", msg.ClientId)
		conn.Close()
		return
	}

	session := &SessionConn{
		SessionId:  uuid.New().String(),
		Client:     client,
		PipModules: msg.PipModules,
		CreatedAt:  time.Now(),
	}
	sessionHub.Register(session)

	//  TODOs
	// search the client which match the requirement and assign the job to that client.
	// send the job execute command to that client.

}
