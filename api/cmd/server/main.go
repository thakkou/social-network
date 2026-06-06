package main

import (
	"log"
	"net/http"

	"01social/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Api)
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
