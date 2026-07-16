package handlers

import (
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
	dblayer "01social/pkg/models/db_layer"
	"01social/pkg/repository" // Import your new repo package
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
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	userLog, err := utilities.ReadJSONRequest[LoginModel](r)
	if err != nil {
		log.Printf("[LOGIN] Error decoding JSON request body: %v", err)
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if userLog.Identifier == "" || userLog.Password == "" {
		log.Printf("[LOGIN] Validation failed: empty identifier or password received")
		utilities.WriteJSON(w, http.StatusBadRequest, "bad credentials", nil)
		return
	}

	// Create Repository using your global db connection
	userRepo := repository.NewUserRepository(db.Database)

	// REPLACED: Raw query with repository GetByIdentifier
	dbUser, err := userRepo.GetByIdentifier(userLog.Identifier)
	if err != nil {
		log.Printf("[LOGIN] User lookup failed for identifier %q: %v", userLog.Identifier, err)
		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.Password),
		[]byte(userLog.Password),
	); err != nil {
		log.Printf("[LOGIN] Password verification failed for user %q: %v", userLog.Identifier, err)
		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	// REPLACED: Raw delete with repository DeleteSessionsByUserID
	err = userRepo.DeleteSessionsByUserID(dbUser.ID)
	if err != nil {
		log.Printf("[LOGIN] Failed to clear old sessions for user ID %d: %v", dbUser.ID, err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	// Create new session
	sessionID := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)

	session := &repository.Session{
		ID:        sessionID,
		ExpiresAt: expiration,
		UserID:    dbUser.ID,
	}

	// REPLACED: Raw insert with repository CreateSession
	err = userRepo.CreateSession(session)
	if err != nil {
		log.Printf("[LOGIN] Failed to write new session to DB for user ID %d: %v", dbUser.ID, err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	ws.NotifyUser(strconv.Itoa(dbUser.ID), "force_logout", nil)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiration,
	})

	utilities.WriteJSON(w, http.StatusOK, "Login Success", map[string]any{
		"id":        dbUser.ID,
		"firstname": dbUser.Firstname,
		"lastname":  dbUser.Lastname,
		"email":     dbUser.Email,
		"birthdate": dbUser.Birthdate,
		"nickname":  dbUser.Nickname,
		"aboutme":   dbUser.AboutMe,
		"avatar":    dbUser.Avatar,
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
	if err != nil {
		log.Printf("[LOGOUT] Session cookie not found in request: %v", err)
		return
	}

	// Create Repository using your global db connection
	userRepo := repository.NewUserRepository(db.Database)

	// REPLACED: utilities.DeleteSession with repository DeleteSessionByID
	err = userRepo.DeleteSessionByID(cookie.Value)
	if err != nil {
		log.Printf("[LOGOUT] Failed to delete session %s from database: %v", cookie.Value, err)
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
	fmt.Println("start registring")

	if r.URL.Path != "/api/register" {
		utilities.WriteJSON(w, http.StatusNotFound, "path not found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	const maxAvatarSize int64 = 1 << 20 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize)

	err := r.ParseMultipartForm(maxAvatarSize)
	if err != nil {
		log.Printf("[REGISTER] Parsing multipart form failed (likely file size limit exceeded): %v", err)
		utilities.WriteJSON(w, http.StatusBadRequest, "Max max size is 1Mb.", nil)
		return
	}

	user := dblayer.User{
		Firstname: strings.TrimSpace(r.FormValue("firstname")),
		Lastname:  strings.TrimSpace(r.FormValue("lastname")),
		Email:     strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		Password:  r.FormValue("password"),
		BirthDate: strings.TrimSpace(r.FormValue("birthDate")),
		Nickname:  strings.TrimSpace(r.FormValue("nickname")),
		AboutMe:   r.FormValue("aboutme"),
	}

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
			log.Printf("[REGISTER] Registration validation failed: empty %s field", field.Name)
			utilities.WriteJSON(w, http.StatusBadRequest, field.Name+" is required", nil)
			return
		}
	}

	if !utilities.IsValidName(user.Firstname) || !utilities.IsValidName(user.Lastname) ||
		!utilities.IsValidEmail(user.Email) || !utilities.IsValidPassword(user.Password) ||
		!utilities.IsValidBirthDate(user.BirthDate) {
		log.Printf("[REGISTER] Validation formats check failed for email %q (User details validation failed)", user.Email)
		utilities.WriteJSON(w, http.StatusBadRequest, "validation fields check failed", nil)
		return
	}

	if user.Nickname != "" && !utilities.IsValidName(user.Nickname) {
		log.Printf("[REGISTER] Nickname format validation failed for: %q", user.Nickname)
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid nickname: use only letters (2-50 characters)", nil)
		return
	}

	if user.AboutMe != "" && !utilities.IsValidDescription(user.AboutMe) {
		log.Printf("[REGISTER] AboutMe bio validation failed for user: %q", user.Email)
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid 'about me': use 2048 characters or less", nil)
		return
	}

	// Create Repository using your global db connection
	userRepo := repository.NewUserRepository(db.Database)

	// REPLACED: Both raw EXISTS checks with your repository IsEmailOrNicknameTaken
	emailExists, nicknameExists, err := userRepo.IsEmailOrNicknameTaken(user.Email, user.Nickname)
	if err != nil {
		log.Printf("[REGISTER] Unique constraints check query failed: %v", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	if emailExists {
		log.Printf("[REGISTER] Conflict detected: email %q is already in use", user.Email)
		utilities.WriteJSON(w, http.StatusBadRequest, "email already exists", nil)
		return
	}
	if nicknameExists {
		log.Printf("[REGISTER] Conflict detected: nickname %q is already taken", user.Nickname)
		utilities.WriteJSON(w, http.StatusBadRequest, "username already taken", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[REGISTER] Failed to hash password: %v", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	var avatarPath string
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			log.Printf("[REGISTER] Rejected file type for avatar upload: %q", ext)
			http.Error(w, "invalid file type", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		if err := os.MkdirAll("uploads/avatars", os.ModePerm); err != nil {
			log.Printf("[REGISTER] Failed to create destination directories: %v", err)
			http.Error(w, "could not process upload", http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(filepath.Join("uploads/avatars", filename))
		if err != nil {
			log.Printf("[REGISTER] Failed to create file on filesystem: %v", err)
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("[REGISTER] Failed copy stream to destination file: %v", err)
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}

		avatarPath = "/uploads/avatars/" + filename
	} else if err != http.ErrMissingFile {
		log.Printf("[REGISTER] Non-standard error retrieving uploaded file: %v", err)
	}

	// Create structural model for the repository
	repoUser := &repository.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  string(hashedPassword),
		Birthdate: user.BirthDate,
		Nickname:  user.Nickname,
		AboutMe:   user.AboutMe,
		Avatar:    avatarPath,
		IsPrivate: 0,
	}

	// REPLACED: Raw INSERT query with repository Create
	err = userRepo.Create(repoUser)
	if err != nil {
		log.Printf("[REGISTER] Failed writing new user records to DB: %v", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "registration success", map[string]string{
		"email": user.Email,
	})
}
