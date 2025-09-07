package main

import (
	"log"
	"net/http"

	"github.com/Amit-syntax/distribute_compute/internal/server"
)


func main() {

	// run a websocket server

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		server.HandleClientConn(w, r)
	})

	log.Println("WebSocket server starting on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("Status endpoint: http://localhost:8080/status")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


