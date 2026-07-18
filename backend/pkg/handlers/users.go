package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/utilities"
)

type userSearchResult struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsPrivate int    `json:"is_private"`
}

type userProfileResponse struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsPrivate int    `json:"is_private"`
}

// get UserProfile
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "api" || parts[1] != "users" {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil || id <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	var user userProfileResponse

	query := `
	SELECT
		id,
		nickname,
		firstname,
		lastname,
		avatar,
		is_private
	FROM USERS
	WHERE id = ?
	`

	err = db.Database.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Firstname,
		&user.Lastname,
		&user.Avatar,
		&user.IsPrivate,
	)

	if err == sql.ErrNoRows {
		utilities.WriteJSON(w, http.StatusNotFound, "user not found", nil)
		return
	}

	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch user", nil)
		return
	}
	utilities.WriteJSON(w, http.StatusOK, "user data fetched", user)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		utilities.WriteJSON(w, http.StatusOK, "users fetched", []userSearchResult{})
		return
	}

	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := db.Database.Query(`
		SELECT id, firstname, lastname, nickname, avatar, is_private
		FROM USERS
		WHERE LOWER(firstname) LIKE ?
		   OR LOWER(lastname) LIKE ?
		   OR LOWER(nickname) LIKE ?
		ORDER BY CASE WHEN LOWER(nickname) = ? THEN 0 ELSE 1 END, firstname COLLATE NOCASE ASC, lastname COLLATE NOCASE ASC
		LIMIT 20
	`, searchTerm, searchTerm, searchTerm, strings.ToLower(query))
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not search users", nil)
		return
	}
	defer rows.Close()

	users := make([]userSearchResult, 0)
	for rows.Next() {
		var u userSearchResult
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Nickname, &u.Avatar, &u.IsPrivate); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not scan users", nil)
			return
		}
		users = append(users, u)
	}

	utilities.WriteJSON(w, http.StatusOK, "users fetched", users)
}

func GetUsernameByToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var userID int
	err = db.Database.QueryRow("SELECT user_id FROM SESSIONS WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not validate session", nil)
		return
	}

	_, err = db.Database.Exec("UPDATE USERS SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not update last seen", nil)
		return
	}

	var (
		firstname sql.NullString
		lastname  sql.NullString
		nickname  sql.NullString
		aboutme   sql.NullString
		avatar    sql.NullString
		lastSeen  sql.NullString
		isPrivate int
	)

	err = db.Database.QueryRow(`
		SELECT firstname, lastname, nickname, aboutme, avatar, is_private, last_seen
		FROM USERS
		WHERE id = ?
	`, userID).Scan(&firstname, &lastname, &nickname, &aboutme, &avatar, &isPrivate, &lastSeen)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not load profile", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "success", map[string]any{
		"authenticated": true,
		"id":            userID,
		"firstname":     firstname.String,
		"lastname":      lastname.String,
		"nickname":      nickname.String,
		"aboutme":       aboutme.String,
		"avatar":        avatar.String,
		"is_private":    isPrivate,
		"last_seen":     lastSeen.String,
	})
}
