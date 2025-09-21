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
		fmt.Printf("Client: %s, IP: %s, Type: %s\n", client.Username, client.IP, client.JoineeType)
	}

}
