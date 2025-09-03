package main

import (
	"log"
	"net/http"

	"github.com/Amit-syntax/distribute_compute/internal/server"
)


func main() {

	// run a websocket server

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWebsocketConn(w, r)
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


