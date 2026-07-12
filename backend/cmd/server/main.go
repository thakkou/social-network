package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	database "01social/pkg/db"
	"01social/pkg/handlers"
	"01social/pkg/routes"
	"01social/pkg/utilities"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, 200, "server is healty", nil)
	fmt.Println("healt")
	return
}

func maxBodySizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const maxSize = 4 * 1024 * 1024
		if r.ContentLength > maxSize {

			utilities.WriteJSON(w, http.StatusRequestEntityTooLarge, "Request body too large", nil)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow common methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Check for refresh command
	refresh := len(os.Args) > 1 && (os.Args[1] == "refresh" || os.Args[1] == "-r")

	if err := database.Init(refresh); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", handlers.HandlerWs)
	http.HandleFunc("/ws/test", handlers.TestBroadcast)
	http.HandleFunc("/assets/", handlers.Static)
	http.HandleFunc("/uploads/", handlers.Static)

	routes.RegisterRoutes()

	log.Println("Server running on http://localhost:8080")

	handler := maxBodySizeMiddleware(http.DefaultServeMux)
	handler = corsMiddleware(handler)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
