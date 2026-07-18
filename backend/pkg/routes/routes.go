package routes

import (
	"net/http"
	"time"

	"01social/pkg/handlers"
	"01social/pkg/middlewares"
)

func RegisterRoutes() {
	// authentification

	http.HandleFunc(
		"/api/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/api/logout",
		middlewares.CheckSessionCookie(handlers.Logout, true),
	)

	http.HandleFunc(
		"/api/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/api/me",
		middlewares.CheckSessionCookie(handlers.GetUsernameByToken, true),
	)

	// profiles
	http.HandleFunc("/api/profile/", middlewares.CheckSessionCookie(handlers.GetProfile, true))
	http.HandleFunc("/api/profile/privacy", middlewares.CheckSessionCookie(handlers.UpdateProfilePrivacy, true))
	// follow
	http.HandleFunc("/api/follow/", middlewares.CheckSessionCookie(handlers.FollowResolver, true))

	// auth providers

	// http.HandleFunc(
	// 	"/api/auth/{provider}",
	// 	middlewares.CheckSessionCookie(handlers.OAuthLogin, false),
	// )

	// http.HandleFunc(
	// 	"/api/auth/{provider}/callback",
	// 	middlewares.CheckSessionCookie(handlers.OAuthCallback, false),
	// )

	// postes

	http.HandleFunc(
		"/api/posts",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetPosts, true),
			3*time.Second,
		),
	)

	// http.HandleFunc(
	// 	"/api/posts/{id}",
	// 	middlewares.RateLimit(
	// 		middlewares.CheckSessionCookie(handlers.GetPostById, true),

	// 		3*time.Second,
	// 	),
	// )

	http.HandleFunc(
		"/api/posts/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreatePost, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/categories",
		middlewares.CheckSessionCookie(handlers.GetAllCategories, true),
	)

	// this resolver for like, dislike, delete, and single post fetch
	http.HandleFunc(
		"/api/posts/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.PostResolver, true),
			250*time.Millisecond),
	)

	// comments

	http.HandleFunc(
		"/api/comments/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateComment, true),
			250*time.Millisecond,
		),
	)

	// this resolver for liking, disliking, delete comment
	http.HandleFunc(
		"/api/comments/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CommentResolver, true),
			500*time.Millisecond),
	)

	// groups
	http.HandleFunc(
		"/api/groups",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.ListGroups, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/groups/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateGroup, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/groups/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GroupResolver, true),
			500*time.Millisecond,
		),
	)

	// user routes
	http.HandleFunc(
		"/api/users/search",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.SearchUsers, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/users/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetUsersById, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/conversations",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetConversation, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/conversation/",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.GetConversationByID, true),
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/messages",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.SendMessage, true),
			100*time.Millisecond,
		),
	)

	http.HandleFunc(
		"/ws",
		middlewares.CheckSessionCookie(handlers.HandlerWs, true),
	)

	http.HandleFunc(
		"/api/online-users",
		middlewares.CheckSessionCookie(handlers.GetOnlineUsers, true),
	)

	http.HandleFunc(
		"/api/notifications",
		middlewares.CheckSessionCookie(handlers.GetNotifications, true),
	)

	http.HandleFunc(
		"/api/notifications/read",
		middlewares.CheckSessionCookie(handlers.MarkNotificationRead, true),
	)

	// http.HandleFunc(
	// 	"/api/users/{id}",
	// 	middlewares.RateLimit(
	// 		middlewares.CheckSessionCookie(handlers.GetPostById, true),
	// 		3*time.Second,
	// 	),
	// )
	// conversation and message conversation
	// http.HandleFunc(
	// 	"/api/messages",
	// 	middlewares.RateLimit(
	// 		middlewares.CheckSessionCookie(handlers.SendMessage, true),
	// 		100*time.Millisecond,
	// 	),
	// )
}
