package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Amit-syntax/distribute_compute/internal/server"
)

func main() {

	// run a websocket server

	httpServer := &http.Server{
		Addr: ":9090",
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		server.HandleClientConn(w, r)
	})

	log.Println("WebSocket server starting on :9090")
	log.Println("WebSocket endpoint: ws://localhost:9090/ws")
	log.Println("Status endpoint: http://localhost:9090/status")

	sysSigCh := make(chan os.Signal, 1)
	signal.Notify(sysSigCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		go func() {
			if err := httpServer.ListenAndServe(); err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()

		select {
		case <-ctx.Done():
			httpServer.Close()
			return
		}

	}()

	go server.CmdLine(ctx, cancel)

	select {
	case <-sysSigCh:
		log.Println("Received interrupt signal, shutting down...")
		cancel()
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	time.Sleep(time.Second * 1)
}
