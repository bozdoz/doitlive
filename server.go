package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer() {
	fmt.Println("[doitlive]", version)
	fmt.Println("[doitlive]", "JS url:", fmt.Sprintf("http://%s/doitlive.js", wsHost))
	fmt.Println("[doitlive]", "WS Endpoint:", fmt.Sprintf("http://%s%s", wsHost, wsEndpoint))

	server := &http.Server{Addr: fmt.Sprintf(":%s", wsPort)}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Unexpected server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for Ctrl-C
	<-stop

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}
}
