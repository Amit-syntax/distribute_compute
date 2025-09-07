package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"regexp"

	"github.com/gorilla/websocket"
)

func Connect(serverIP string, port string) {
	url := fmt.Sprintf("ws://%s:%s/register", serverIP, port)

	// connect to websocket server.
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)

	if err != nil {
		log.Printf("Failed to connect to server: %v", err)
		return  
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

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// holding cli to read user commands
	go readUserCommands(wg)

	wg.Wait()
}


func readUserCommands(wg *sync.WaitGroup) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("cmd: ")
		command, _ := reader.ReadString('\n')
		log.Printf("executing command: %v", command)

		if command == "exit\n" {
			wg.Done()
			return 
		}
	}
}


func isValidUsername(username string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
    if !re.MatchString(username) {
        return false
    }
    if containsSpace := regexp.MustCompile(`\s`).MatchString(username); containsSpace {
        return false
    }
    return true
}

func sendRegisterMsg(conn *websocket.Conn) error {
	registerRole := ""
	clientUsername := ""

	// validations for role and username
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter role (worker/consumer): ")
		registerRole, _ = reader.ReadString('\n')
		fmt.Print("Enter name: ")
		clientUsername, _ = reader.ReadString('\n')

		if strings.Trim(registerRole, "\n") == "worker" || strings.Trim(registerRole, "\n") == "consumer" {
			registerRole = strings.Trim(registerRole, "\n")
			clientUsername = strings.Trim(clientUsername, "\n")
			break
		} else {
			fmt.Println("Invalid role. Please enter 'worker' or 'consumer'.")
		}

		if !isValidUsername(clientUsername) {
			fmt.Println("Invalid username. Please use 3-20 characters: letters, numbers, and underscores only.")
			continue
		}
	}

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

	defer func () {
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {

			log.Printf("Error reading message: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v", err)
			}

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


