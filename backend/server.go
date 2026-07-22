package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "01social/pkg/db/sqlite"
	"01social/pkg/handlers"
	"01social/pkg/repository"
	"01social/pkg/routes"
	"01social/pkg/utilities"
)

// @title Social Network API
// @version 1.0
// @description A social network API built with Go. Supports user authentication, posts, comments, reactions, groups, messaging, notifications, and real-time WebSocket events.
// @termsOfService https://example.com/terms
//
// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey SessionCookie
// @in cookie
// @name session_id
//
// @tag.name Authentication
// @tag.description Login, Register, Logout, Session validation
//
// @tag.name Posts
// @tag.description Create, read, filter, and manage posts
//
// @tag.name Comments
// @tag.description Create and manage comments on posts
//
// @tag.name Reactions
// @tag.description Like/dislike posts and comments
//
// @tag.name Follow
// @tag.description Follow/unfollow users and manage follow requests
//
// @tag.name Profile
// @tag.description View and update user profiles
//
// @tag.name Groups
// @tag.description Create and manage groups, group posts, events, invites
//
// @tag.name Conversations
// @tag.description Direct messaging and group chat
//
// @tag.name Notifications
// @tag.description View and manage notifications
//
// @tag.name Search
// @tag.description Search users and groups
//
// @tag.name Categories
// @tag.description Post categories
//
// @tag.name WebSocket
// @tag.description Real-time events via WebSocket connections

// healthHandler responds with the server health status.
// @Summary Health check
// @Description Returns the server health status.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Server is healthy"
// @Router /health [get]
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

	if err := db.Init(refresh); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	// here init reposotory
	repos := repository.NewRepositories(db.Database)
	handlers.Init(repos)
	http.HandleFunc("/health", healthHandler)
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
