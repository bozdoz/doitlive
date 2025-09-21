package main

import (
	_ "embed"
	"flag"

	"fmt"
	"net/http"
)

//go:embed version.txt
var version string

var (
	wsPort      = "35729"
	wsHost      = fmt.Sprintf("localhost:%s", wsPort)
	wsEndpoint  = "/ws/reload"
	broadcast   = make(chan string)
	HardRefresh = "false"
	ignored     = map[string]struct{}{
		".git":         {},
		"node_modules": {},
	}
	is_hard_refresh = flag.Bool("hard-refresh", false, "JS code executes a hard refresh")
	is_debug        = flag.Bool("debug", false, "set debug")
)

func debug(a ...any) {
	a = append([]any{"[doitlive]"}, a...)
	if *is_debug {
		fmt.Println(a...)
	}
}

func main() {
	flag.Parse()

	if *is_hard_refresh {
		HardRefresh = "true"
	}

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
