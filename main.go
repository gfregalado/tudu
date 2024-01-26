package main

import (
	"context"
	"errors"
	"github.com/gfregalado/todo/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Set up a channel to listen for interrupt or terminate signals from the OS.
	// This allows for graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Create a server and set timeout values for more robust handling of connections.
	server := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      nil, // Use http.DefaultServeMux
	}

	// Create a file server for static files
	fs := http.FileServer(http.Dir("static"))
	// Serve static files on /static/ route
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		http.HandleFunc("/", handlers.HomeGetHandler)

		log.Println("Server is ready to handle requests at http://localhost:3000")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", ":3000", err)
		}
	}()

	// Listen for an interrupt signal from the OS.
	<-stopChan
	log.Println("Server is shutting down...")

	// Give outstanding requests a deadline for completion.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop the HTTP server.
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("There was an error shutting down: %v", err)
	}

	log.Println("Server stopped")
}
