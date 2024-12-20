package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BarkinBalci/golangassignment/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable is not set")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/mongo", handlers.MongoHandler)
	http.HandleFunc("/memory", handlers.MemoryHandler)

	log.Printf("Server started on port: %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Fatalf("error starting server %v", err)
	}
}
