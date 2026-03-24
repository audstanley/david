package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/audstanley/david/app/api"
)

func main() {
	// Load config from environment
	port := os.Getenv("DAVID_PORT")
	if port == "" {
		port = "8080"
	}

	// Create router
	router := api.NewRouter(api.Config{
		Port: port,
	})
	router.SetupRoutes()

	// Start server
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting david API server on %s\n", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
