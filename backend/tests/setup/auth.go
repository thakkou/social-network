package setup

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// CreateSession inserts a new session for the given user into the database
// and returns the session ID. The session expires in 24 hours.
func CreateSession(db *sql.DB, userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := db.Exec(
		"INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)",
		sessionID,
		expiresAt.Format("2006-01-02 15:04:05"),
		userID,
	)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// AuthCookie creates an http.Cookie with the given session ID.
// Use this to attach authentication to test requests.
func AuthCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
}

// AddAuthCookie adds the session cookie to the given request.
func AddAuthCookie(r *http.Request, sessionID string) {
	r.AddCookie(AuthCookie(sessionID))
}

// DeleteSession removes a session from the database.
func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM SESSIONS WHERE id = ?", sessionID)
	return err
}
