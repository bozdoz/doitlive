package main

import (
	_ "embed"
	"sync"

	"fmt"
	"net/http"

	// web sockets
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = make(map[*websocket.Conn]struct{})
	mu       sync.Mutex
)

func handleWebSocketConnect(w http.ResponseWriter, r *http.Request) {
	// allow all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	fmt.Println("[doitlive]", "Client Connected")

	mu.Lock()
	clients[ws] = struct{}{}
	mu.Unlock()

	for {
		// wait for client message to remove them
		if _, _, err := ws.ReadMessage(); err != nil {
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			fmt.Println("[doitlive]", "Client Disconnected")
			break
		}
	}
}

func waitForChanges() {
	for {
		// wait for broadcast
		msg := <-broadcast
		mu.Lock()
		// send to clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("WebSocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
