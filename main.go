package main

import (
	_ "embed"

	"fmt"
	"net/http"
)

//go:embed version.txt
var version string

var (
	wsPort     = "35729"
	wsHost     = fmt.Sprintf("localhost:%s", wsPort)
	wsEndpoint = "/ws/reload"
	broadcast  = make(chan string)
	ignored    = map[string]struct{}{
		".git":         {},
		"node_modules": {},
	}
)

func main() {
	// WebSocket endpoint (see ./wshandler.go)
	http.HandleFunc(wsEndpoint, handleWebSocketConnect)

	// Serve the client JS (see ./jshandler.go)
	http.HandleFunc("/doitlive.js", handleJS)

	// Start broadcaster
	go waitForChanges()
	// start file watcher (see ./filewatcher.go)
	go watchFiles()

	// start server
	startServer()

	// remove dangling "%" in the terminal
	fmt.Println("")
}
