package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func Connect(serverIP string, port string) {
	url := fmt.Sprintf("ws://%s:%s", serverIP, port)

	// connect to websocket server.
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)

	if err != nil {
		log.Printf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server: %s", url)

	err = sendRegisterMsg(conn)
	if err != nil {
		log.Printf("Failed to send register message: %v", err)
		return
	}

	// ask for name
	go recvMsg(conn)

}


func sendRegisterMsg(conn *websocket.Conn) error {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter role (worker/consumer): ")
	registerRole, _ := reader.ReadString('\n')
	fmt.Print("Enter name: ")
	clientUsername, _ := reader.ReadString('\n')

	registerMsg := map[string]any{
		"action": registerRole,
		"client_username":  clientUsername,
		"ip":               "", //TODO: get this
		"joinee_role":      registerRole,
	}

	regMsgByte, err := json.Marshal(registerMsg)
	err = conn.WriteMessage(websocket.TextMessage, regMsgByte)
	if err != nil {
		return err
	}
	log.Printf("Sent register message: %s", registerMsg)

	return nil
}


func recvMsg(conn *websocket.Conn) {

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		log.Printf("Received message: %s", message)
		err = executeCmd(message)
		if err != nil {
			log.Printf("Error executing command: %v", err)
		}
	}
}

func executeCmd(cmd []byte) error {

	jsonCmd := make(map[string]any)
	err := json.Unmarshal(cmd, &jsonCmd)
	if err != nil {
		log.Printf("Error unmarshalling command: %v", err)
		return err
	}

	log.Printf("Executing command: %s", cmd)

	return nil

}


