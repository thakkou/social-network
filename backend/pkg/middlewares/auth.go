package middlewares

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	db "01social/pkg/db/sqlite"
	"01social/pkg/utilities"
)

func CheckSessionCookie(handler http.HandlerFunc, requiresAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		for _, cookie := range r.Cookies() {
			log.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
		}
		if err != nil || cookie.Value == "" {
			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "login required", nil)
				return
			}
			handler(w, r)
			return
		}

		var userId int
		var expiryTime time.Time
		err = db.Database.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&userId, &expiryTime)

		switch err {
		case nil:
			// continue below
		case sql.ErrNoRows:
			utilities.ClearSessionCookie(w)
			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "login required", nil)
				return
			}
			handler(w, r)
			return
		default:
			utilities.WriteJSON(w, http.StatusInternalServerError, "database error", nil)
			return
		}

		if expiryTime.Before(time.Now()) {
			utilities.DeleteSession(cookie.Value)
			utilities.ClearSessionCookie(w)
			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "session expired", nil)
				return
			}
			handler(w, r)
			return
		}

		// Valid session
		if requiresAuth {
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			handler(w, r.WithContext(ctx)) // <-- inject userID here
			return
		}

		utilities.WriteJSON(w, http.StatusConflict, "Unauthorized", nil)
	}
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}
