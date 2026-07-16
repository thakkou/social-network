package repository

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Birthdate  string    `json:"birthdate"`
	Nickname   string    `json:"nickname"`
	AboutMe    string    `json:"aboutme"`
	Avatar     string    `json:"avatar"`
	IsPrivate  int       `json:"is_private"`
	LastSeen   time.Time `json:"last_seen"`
	is_private int
}

type Session struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"user_id"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// here because the thire are optionale
func (r *UserRepository) Create(u *User) error {
	nickname := sql.NullString{String: u.Nickname, Valid: u.Nickname != ""}
	aboutme := sql.NullString{String: u.AboutMe, Valid: u.AboutMe != ""}
	avatar := sql.NullString{String: u.Avatar, Valid: u.Avatar != ""}

	query := `INSERT INTO USERS (firstname, lastname, email, password, birthdate, nickname, aboutme, avatar) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, u.Firstname, u.Lastname, u.Email, u.Password, u.Birthdate, nickname, aboutme, avatar, u.IsPrivate)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}
	return nil
}

func (r *UserRepository) GetByIdentifier(identifier string) (*User, error) {
	query := `SELECT id, firstname, lastname, email, password, birthdate, nickname, aboutme, avatar, is_private FROM USERS WHERE email = ? OR nickname = ?`
	var u User
	var nickname, aboutme, avatar sql.NullString
	err := r.DB.QueryRow(query, identifier, identifier).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Password, &u.Birthdate, &nickname, &aboutme, &avatar, &u.IsPrivate)
	if err != nil {
		return nil, err
	}
	u.Nickname = nickname.String
	u.AboutMe = aboutme.String
	u.Avatar = avatar.String
	return &u, nil
}

func (r *UserRepository) CreateSession(s *Session) error {
	query := `INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, s.ID, s.ExpiresAt.Format("2006-01-02 15:04:05"), s.UserID)
	return err
}

func (r *UserRepository) DeleteSessionsByUserID(userID int) error {
	_, err := r.DB.Exec("DELETE FROM SESSIONS WHERE user_id = ?", userID)
	return err
}

func (r *UserRepository) DeleteSessionByID(sessionID string) error {
	_, err := r.DB.Exec("DELETE FROM SESSIONS WHERE id = ?", sessionID)
	return err
}

func (r *UserRepository) IsEmailOrNicknameTaken(email, nickname string) (bool, bool, error) {
	var emailExists, nicknameExists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE email = ?)", email).Scan(&emailExists)
	if err != nil {
		return false, false, err
	}
	if nickname != "" {
		err = r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE nickname = ? COLLATE NOCASE)", nickname).Scan(&nicknameExists)
		if err != nil {
			return emailExists, false, err
		}
	}
	return emailExists, nicknameExists, nil
}

func (r *UserRepository) UpdatePrivacy(userID int, isPrivate int) error {
	_, err := r.DB.Exec("UPDATE USERS SET is_private = ? WHERE id = ?", isPrivate, userID)
	return err
}
