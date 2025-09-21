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
	is_hard_refresh = flag.Bool("hard-refresh", false, "JS code executes a hard refresh")
	is_debug        = flag.Bool("debug", false, "set debug")
	proxy_port      = flag.Int("proxy", 4000, "proxy port")
	host_port       = flag.Int("host", 8000, "host port")
)

var (
	HardRefresh = false
	wsEndpoint  = "/ws/reload"
	broadcast   = make(chan string)
	ignored     = map[string]struct{}{
		".git":         {},
		"node_modules": {},
	}
)

func init() {
	flag.Parse()
}

func debug(a ...any) {
	a = append([]any{"[doitlive]"}, a...)
	if *is_debug {
		fmt.Println(a...)
	}
}

func main() {
	http.Handle("/", runProxy())

	// WebSocket endpoint (see ./wshandler.go)
	http.HandleFunc(wsEndpoint, handleWebSocketConnect)

	// Start broadcaster
	go waitForChanges()
	// start file watcher (see ./filewatcher.go)
	go watchFiles()

	fmt.Println("[doitlive]", version)
	fmt.Println("[doitlive]", "Host:", fmt.Sprintf("http://localhost:%d", *host_port))
	fmt.Println("[doitlive]", "Proxied:", fmt.Sprintf("http://localhost:%d", *proxy_port))

	// start server
	startServer()

	// remove dangling "%" in the terminal
	fmt.Println("")
}
