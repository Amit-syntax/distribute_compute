package server

import (
	"context"
	"fmt"
)

func CmdLine(ctx context.Context, cancel context.CancelFunc) {

	for {
		// Read input from the command line.
		var command string
		fmt.Print("> ")
		if _, err := fmt.Scanln(&command); err != nil {
			// Handle potential errors, e.g., EOF
			if err.Error() == "unexpected newline" {
				continue
			}
			fmt.Println("CLI input error:", err)
			return
		}

		if command == "exit" {
			cancel()
			return
		}

		if command == "clients_list" {
			fmt.Print("Listing connected clients: \n")
			listClients()
		}

		if command == "sessions_list" {
			fmt.Print("Listing active consumer sessions: \n")
			listSessions()
		}

		if command == "help" {
			fmt.Println("Available commands:")
			fmt.Println("  exit            - Shut down the server")
			fmt.Println("  clients_list    - List all connected clients")
			fmt.Println("  sessions_list   - List all active consumer sessions")
			fmt.Println("  help            - Show this help message")
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func listClients() {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if len(hub.clients) == 0 {
		fmt.Println("No connected clients.")
		return
	}

	for client := range hub.clients {
		fmt.Printf("Client: %s, Type: %s, Conn: %v\n", client.Username, client.JoineeType, hub.clients[client])
	}

}

func listSessions() {
	consumerSessionHub.mu.RLock()
	defer consumerSessionHub.mu.RUnlock()

	if len(consumerSessionHub.sessions) == 0 {
		fmt.Println("No active consumer sessions.")
		return
	}
	for sessionId, session := range consumerSessionHub.sessions {
		fmt.Printf("Session ID: %s, Client: %s, CreatedAt: %v\n", sessionId, session.Client.Username, session.CreatedAt)
	}
}
