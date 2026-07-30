package setup

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"time"

	"01social/pkg/db/sqlite"
	"01social/pkg/handlers"
	"01social/pkg/middlewares"
	"01social/pkg/repository"
)

// NewTestServer creates a fully configured test HTTP server.
// It:
// 1. Injects the test database into the global sqlite package
// 2. Initializes repositories with the test DB
// 3. Initializes handlers with repositories
// 4. Registers all routes
// 5. Returns an httptest.Server ready for requests
//
// Call server.Close() when done to release resources.
func NewTestServer(db *sql.DB) *httptest.Server {
	// Inject the test database so handlers using sqlite.DB() work correctly
	sqlite.SetDB(db)

	// Initialize repositories
	repos := repository.NewRepositories(db)
	handlers.Init(repos)

	// Use a custom mux so test routes don't pollute the DefaultServeMux
	mux := http.NewServeMux()
	registerTestRoutes(mux, repos)

	// Create the test server
	server := httptest.NewServer(mux)
	return server
}

// registerTestRoutes registers all application routes on the given mux.
// This mirrors routes.RegisterRoutes() but uses a custom ServeMux.
func registerTestRoutes(mux *http.ServeMux, repos *repository.Repositories) {
	// Rate limiting is effectively disabled in tests by using 1 ns intervals.
	// This allows tests to run at full speed while still exercising the middleware.
	rateLimitDisabled := time.Nanosecond

	// Authentification
	mux.HandleFunc("/api/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/logout",
		middlewares.CheckSessionCookie(handlers.Logout, true),
	)
	mux.HandleFunc("/api/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/me",
		middlewares.CheckSessionCookie(handlers.GetUsernameByToken, true),
	)
	mux.HandleFunc("/api/session/validate",
		middlewares.CheckSessionCookie(handlers.ValidateSession, false),
	)

	// Search
	mux.HandleFunc("/api/search",
		middlewares.CheckSessionCookie(handlers.FindQuery, true),
	)

	// Profiles
	mux.HandleFunc("/api/profile/update",
		middlewares.CheckSessionCookie(handlers.UpdateProfile, true),
	)
	mux.HandleFunc("/api/profile/",
		middlewares.CheckSessionCookie(handlers.GetProfile, true),
	)
	mux.HandleFunc("/api/profile/privacy",
		middlewares.CheckSessionCookie(handlers.UpdateProfilePrivacy, true),
	)

	// Follow
	mux.HandleFunc("/api/follow/",
		middlewares.CheckSessionCookie(handlers.FollowResolver, true),
	)

	// Posts
	mux.HandleFunc("/api/posts",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetPosts, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/posts/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreatePost, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/categories",
		middlewares.CheckSessionCookie(handlers.GetAllCategories, true),
	)
	mux.HandleFunc("/api/posts/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.PostResolver, true),
			rateLimitDisabled,
		),
	)

	// Comments
	mux.HandleFunc("/api/comments/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateComment, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/comments/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CommentResolver, true),
			rateLimitDisabled,
		),
	)

	// Groups
	mux.HandleFunc("/api/groups",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.ListGroups, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateGroup, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GroupResolver, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/public/{id}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetGroupPublic, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/content/{id}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetGroupContent, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/members/{id}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetGroupMembers, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/groups/invite-candidates/{id}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetInviteCandidates, true),
			rateLimitDisabled,
		),
	)

	// Users
	mux.HandleFunc("/api/users/search",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.SearchUsers, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/users/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetUsersById, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/users/groups",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetMyGroups, true),
			rateLimitDisabled,
		),
	)

	// Conversations & Messages
	mux.HandleFunc("/api/conversations",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetConversation, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/conversation/{type}/{id}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetConversationByID, true),
			rateLimitDisabled,
		),
	)
	mux.HandleFunc("/api/messages",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.SendMessage, true),
			rateLimitDisabled,
		),
	)

	// WebSocket
	mux.HandleFunc("/api/ws-ticket",
		middlewares.CheckSessionCookie(handlers.CreateWsTicket, true),
	)
	mux.HandleFunc("/ws", handlers.HandlerWs)

	// Online users
	mux.HandleFunc("/api/online-users",
		middlewares.CheckSessionCookie(handlers.GetOnlineUsers, true),
	)

	// Notifications
	mux.HandleFunc("/api/notifications",
		middlewares.CheckSessionCookie(handlers.GetNotifications, true),
	)
	mux.HandleFunc("/api/notifications/readAll",
		middlewares.CheckSessionCookie(handlers.MarkAllNotificationsRead, true),
	)
	mux.HandleFunc("/api/notifications/read",
		middlewares.CheckSessionCookie(handlers.MarkNotificationRead, true),
	)
	mux.HandleFunc("/api/notifications/delete",
		middlewares.CheckSessionCookie(handlers.DeletNotif, true),
	)
	mux.HandleFunc("/api/notifications/deletAll",
		middlewares.CheckSessionCookie(handlers.DeletAllNotif, true),
	)

	// Documentation
	mux.HandleFunc("/api/docs/swagger.json", handlers.GetSwaggerJSON)
	mux.HandleFunc("/api/docs", handlers.GetSwaggerUI)
}

// NewTestServerWithDefaults creates a test server with a fresh in-memory DB
// that has been migrated and seeded. This is the simplest way to get a
// fully functional test server.
func NewTestServerWithDefaults() (*httptest.Server, *sql.DB, error) {
	db, err := NewTestDB()
	if err != nil {
		return nil, nil, err
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, nil, err
	}

	if err := SeedTestData(db); err != nil {
		db.Close()
		return nil, nil, err
	}

	server := NewTestServer(db)
	return server, db, nil
}
