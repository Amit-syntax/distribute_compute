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

type ConsumerSessionMsg struct {
	PipModules []string `json:"pip_modules"`
	ClientId   string   `json:"client_id"`
}

type ConsumerSessionConn struct {
	SessionId   string
	SessionConn *websocket.Conn
	Client      *Client
	PipModules  []string
	CreatedAt   time.Time
}

type ConsumerSessionHub struct {
	sessions map[string]*ConsumerSessionConn
	mu       *sync.RWMutex
}

func NewConsumerSessionHub() *ConsumerSessionHub {
	return &ConsumerSessionHub{
		sessions: make(map[string]*ConsumerSessionConn),
		mu:       &sync.RWMutex{},
	}
}

func (h *ConsumerSessionHub) Register(session *ConsumerSessionConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sessions[session.SessionId] = session
}

var consumerSessionHub = NewConsumerSessionHub()

func SessionHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Wait for client registration message
	conn.SetReadLimit(512)

	_, msgByte, err := conn.ReadMessage()
	log.Print("SessionHandler received message: ", string(msgByte))
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}

	msg := &ConsumerSessionMsg{}
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

	session := &ConsumerSessionConn{
		SessionId:   uuid.New().String(),
		SessionConn: conn,
		Client:      client,
		PipModules:  msg.PipModules,
		CreatedAt:   time.Now(),
	}
	consumerSessionHub.Register(session)

	sessionMsg := common.Message{
		Type:        common.SessionAckMsgType,
		Description: "Session created",
		Body:        common.SessionAckMsg{SessionId: session.SessionId},
	}
	log.Print("Session created with ID: ", session.SessionId)
	// send session ID to client in ack msg
	if err := conn.WriteJSON(sessionMsg); err != nil {
		log.Printf("error sending ack: %v", err)
		return
	}

	go session.readSession()
}

func (s *ConsumerSessionConn) readSession() {

	defer func() {
		s.Client.hub.Unregister(s.Client)
		s.SessionConn.Close()
	}()

	for {
		_, msgByte, err := s.SessionConn.ReadMessage()
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

		if msg.Type != common.SessionRemoteExecReq {
			log.Printf("unknown message type: %v", msg.Type)
			continue
		}
		log.Printf("Received remote exec request: %v", msg)

		executionUUID := sendRemoteExecRequestToWorker(s, msg)

		respMsg := &common.Message{
			Type:        common.SessionRemoteExecInit,
			Description: "Remote execution request initialized",
			Body: common.SessionRemoteExecInitMsg{
				ExecutionID: executionUUID,
			},
		}
		if err := s.SessionConn.WriteJSON(respMsg); err != nil {
			log.Printf("error sending remote exec init response: %v", err)
			continue
		}
		log.Printf("Sent remote exec init response with execution ID: %s", executionUUID)

	}
}

//  TODOs
// search the client which match the requirement and assign the job to that client.
// send the job execute command to that client.

func sendRemoteExecRequestToWorker(session *ConsumerSessionConn, msg *common.Message) string {
	log.Printf("Sending remote exec request to client: %v", session.Client)

	executionUUID := uuid.New().String()

	// TODO: Find a worker client
	// TODO: sending the request to the client.

	return executionUUID
}
