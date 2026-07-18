package repository

import (
	"database/sql"
	"log"
	"strings"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// GetPublicProfile returns only the fields visible to everyone.
func (r *ProfileRepository) GetPublicProfile(userID int) (*User, error) {
	var u User

	query := `
        SELECT
            id,
            firstname,
            lastname,
            nickname,
            avatar,
            is_private
        FROM USERS
        WHERE id = ?
    `

	var nickname, avatar sql.NullString

	err := r.DB.QueryRow(query, userID).Scan(
		&u.ID,
		&u.Firstname,
		&u.Lastname,
		&nickname,
		&avatar,
		&u.IsPrivate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[PROFILE] user %d not found", userID)
		}
		return nil, err
	}

	u.Nickname = nickname.String
	u.Avatar = avatar.String

	return &u, nil
}

// GetProfile returns the complete profile.
func (r *ProfileRepository) GetProfile(userID int) (*User, error) {
	var u User
	query := `
        SELECT
            id,
            created_at,
            firstname,
            lastname,
            email,
            birthdate,
            nickname,
            aboutme,
            avatar,
            is_private
        FROM USERS
        WHERE id = ?
    `

	var nickname, aboutme, avatar sql.NullString

	err := r.DB.QueryRow(query, userID).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.Firstname,
		&u.Lastname,
		&u.Email,
		&u.Birthdate,
		&nickname,
		&aboutme,
		&avatar,
		&u.IsPrivate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[PROFILE] user %d not found", userID)
		}
		return nil, err
	}

	u.Nickname = nickname.String
	u.AboutMe = aboutme.String
	u.Avatar = avatar.String

	return &u, nil
}

// SearchProfiles searches users by first name, last name, nickname (username),
// or full name.
func (r *ProfileRepository) SearchProfiles(text string) ([]User, error) {
	query := `
	SELECT
		id,
		firstname,
		lastname,
		nickname,
		avatar,
		is_private
	FROM USERS
	WHERE
		LOWER(firstname) LIKE LOWER(?)
		OR LOWER(lastname) LIKE LOWER(?)
		OR LOWER(COALESCE(nickname, '')) LIKE LOWER(?)
		OR LOWER(firstname || ' ' || lastname) LIKE LOWER(?)
	ORDER BY firstname, lastname
	LIMIT 20
	`

	search := "%" + strings.TrimSpace(text) + "%"

	rows, err := r.DB.Query(query, search, search, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		var nickname, avatar sql.NullString

		if err := rows.Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
			&nickname,
			&avatar,
			&u.IsPrivate,
		); err != nil {
			return nil, err
		}

		u.Nickname = nickname.String
		u.Avatar = avatar.String

		users = append(users, u)
	}

	return users, rows.Err()
}
