package repository

import (
	"database/sql"
	"log"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// GetProfile returns limited or full profile
func (r *ProfileRepository) GetProfile(userID int, greedy bool) (*User, error) {
	var u User

	if greedy {
		query := `
			SELECT id, firstname, lastname, nickname, is_private
			FROM USERS
			WHERE id = ?
		`

		var nickname sql.NullString

		err := r.DB.QueryRow(query, userID).Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
			&nickname,
			&u.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		u.Nickname = nickname.String

		return &u, nil
	}

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
