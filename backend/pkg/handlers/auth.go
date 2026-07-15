package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	db "01social/pkg/db/sqlite"
	"01social/pkg/models"
	"01social/pkg/utilities"
	"01social/pkg/ws"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/login" {
		utilities.WriteJSON(w, http.StatusNotFound, "path not found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
		return
	}

	type LoginModel struct {
		Identifier string `json:"identifier"` // "identifier"
		Password   string `json:"password"`
	}

	userLog, err := utilities.ReadJSONRequest[LoginModel](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if userLog.Identifier == "" || userLog.Password == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "bad credentials", nil)
		return
	}

	var (
		userID                                                           int
		firstname, lastname, email, birthDate, nickname, aboutme, avatar string
		hashedPassword                                                   sql.NullString
	)

	err = db.Database.QueryRow( // + created at
		`SELECT id, firstname, lastname, email, password, birthdate, nickname, aboutme, avatar
		 FROM users
		 WHERE email = ? OR nickname = ?`,
		userLog.Identifier,
		userLog.Identifier,
	).Scan(&userID, &firstname, &lastname, &email, &hashedPassword, &birthDate, &nickname, &aboutme, &avatar)
	if err != nil {
		log.Printf("[LOGIN] User not found: %q (%v)", userLog.Identifier, err)

		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	if !hashedPassword.Valid {
		log.Printf("[LOGIN] User %q has no valid password hash", userLog.Identifier)

		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword.String),
		[]byte(userLog.Password),
	); err != nil {
		log.Printf("[LOGIN] Invalid password for user %q password=%q", userLog.Identifier, userLog.Password)
		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	// Remove old sessions
	_, err = db.Database.Exec(
		"DELETE FROM sessions WHERE user_id = ?",
		userID,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	// Create new session
	sessionID := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = db.Database.Exec(
		"INSERT INTO sessions (id, expires_at, user_id) VALUES (?, ?, ?)",
		sessionID,
		expiration,
		userID,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	ws.NotifyUser(strconv.Itoa(userID), "force_logout", nil)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiration,
	})

	utilities.WriteJSON(w, http.StatusOK, "Login Success", map[string]any{
		"id": userID, // ==> not needed
		// "createdAt":
		"firstname": firstname,
		"lastname":  lastname,
		"email":     email,
		"birthdate": birthDate,
		"nickname":  nickname, // nickname can be empty, so should pass first and last name instead !
		"aboutme":   aboutme,
		"avatar":    avatar,
		"token":     sessionID,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/logout" {
		utilities.WriteJSON(w, 404, `path not found`, nil)
		return
	}
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, `Method not allowed`, nil)

		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		return
	}

	err = utilities.DeleteSession(cookie.Value)
	if err != nil {
		log.Println(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	utilities.WriteJSON(w, 201, `log out succes`, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Check route
	if r.URL.Path != "/api/register" {
		utilities.WriteJSON(w, http.StatusNotFound, "path not found", nil)
		return
	}

	// Check method
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// Content type (optional but fine to keep)
	// if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
	// 	utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
	// 	return
	// }

	// added to handle avatar !
	// err := r.ParseMultipartForm(5 << 20) // 5MB max
	// if err != nil {
	// 	utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }
	///

	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2 MB hardcoded

	user := models.User{
		Firstname: strings.TrimSpace(r.FormValue("firstname")),
		Lastname:  strings.TrimSpace(r.FormValue("lastname")),
		Email:     strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		Password:  r.FormValue("password"),
		BirthDate: strings.TrimSpace(r.FormValue("birthDate")),
		Nickname:  strings.TrimSpace(r.FormValue("nickname")),
		AboutMe:   r.FormValue("aboutme"),
	}

	// ✅ REPLACED PART (clean)
	// user, err := utilities.ReadJSONRequest[models.User](r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	// // Normalize input
	// user.Firstname = strings.TrimSpace(user.Firstname)
	// user.Lastname = strings.TrimSpace(user.Lastname)
	// user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	// user.BirthDate = strings.TrimSpace(user.BirthDate)
	// // -- optional
	// user.Nickname = strings.TrimSpace(user.Nickname)
	// user.AboutMe = strings.TrimSpace(user.AboutMe)

	const maxAvatarSize int64 = 1 << 20 // 1 MB

	err := r.ParseMultipartForm(maxAvatarSize)
	// ParseMultipartForm sets the in-memory buffer limit.
	// If the file exceeds that limit, Go silently spills the overflow to a temp file on disk.
	if err != nil {
		fmt.Println(err)
		utilities.WriteJSON(w, http.StatusBadRequest, "Max max size is 1Mb.", nil)
		return
	}

	// Check required fields
	fields := []struct {
		Name  string
		Value string
	}{
		{"firstname", user.Firstname},
		{"lastname", user.Lastname},
		{"email", user.Email},
		{"password", user.Password},
		{"birthdate", user.BirthDate},
	}

	for _, field := range fields {
		if field.Value == "" {
			utilities.WriteJSON(w, http.StatusBadRequest, field.Name+" is required", nil)
			return
		}
	}

	// Validate fields
	if !utilities.IsValidName(user.Firstname) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid first name: use only letters (2-50 characters)", nil)
		return
	}

	if !utilities.IsValidName(user.Lastname) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid last name: use only letters (2-50 characters)", nil)
		return
	}

	if !utilities.IsValidEmail(user.Email) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid email: must be a valid email address", nil)
		return
	}

	if !utilities.IsValidPassword(user.Password) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid password: must be 6-25 characters", nil)
		return
	}

	if !utilities.IsValidBirthDate(user.BirthDate) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid birth date", nil)
		return
	}

	// add optional
	if user.Nickname != "" && !utilities.IsValidName(user.Nickname) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid nickname: use only letters (2-50 characters)", nil)
		return
	}

	if user.AboutMe != "" && !utilities.IsValidDescription(user.AboutMe) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid 'about me': use 2048 characters or less", nil)
		return
	}

	// Check email exists
	var emailExists bool
	err = db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		user.Email,
	).Scan(&emailExists)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	if emailExists {
		utilities.WriteJSON(w, http.StatusBadRequest, "email already exists", nil)
		return
	}

	// Check nickname exists (add optional)
	if user.Nickname != "" {
		var nicknameExists bool
		err = db.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ? COLLATE NOCASE)",
			user.Nickname,
		).Scan(&nicknameExists)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
			return
		}
		if nicknameExists {
			utilities.WriteJSON(w, http.StatusBadRequest, "username already taken", nil)
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	// store avatar:
	var avatarPath string
	file, header, err := r.FormFile("avatar")
	if err == nil { // avatar is optional — err != nil just means none was sent
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			http.Error(w, "invalid file type", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		if err := os.MkdirAll("uploads/avatars", os.ModePerm); err != nil {
			log.Fatal("failed to create upload directory:", err)
		}
		dst, err := os.Create(filepath.Join("uploads/avatars", filename))
		if err != nil {
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}

		avatarPath = "/uploads/avatars/" + filename
	} else {
		fmt.Println(err)
		fmt.Println("No image uploaded, continuing without it")
	}

	// Insert user
	_, err = db.Database.Exec(
		`INSERT INTO users (firstname, lastname, email, password, birthdate, nickname, aboutme, avatar)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Firstname,
		user.Lastname,
		user.Email,
		string(hashedPassword),
		user.BirthDate,
		user.Nickname,
		user.AboutMe,
		avatarPath,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		fmt.Println(err)
		return
	}
	type RegisterResponse struct { // needed only here !!!?
		// Nickname string `json:"nickname"`
		Email string `json:"email"`
	}
	response := RegisterResponse{
		// Nickname: user.Nickname,
		Email: user.Email,
	}

	utilities.WriteJSON(w, http.StatusOK, "registration success", response)
}
