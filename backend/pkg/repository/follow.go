package repository

import (
	"database/sql"
	"fmt"
)

type Follow struct {
	FollowerID  int    `json:"follower_id"`
	FollowingID int    `json:"following_id"`
	Status      string `json:"status"` // pending | accepted
}

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

// Follow creates a follow relationship.
// For public users: status = accepted
// For private users: status = pending
func (r *FollowRepository) Follow(followerID, followingID int, status string) error {
	query := `
		INSERT INTO FOLLOWS 
		(follower_id, following_id, status)
		VALUES (?, ?, ?)

		ON CONFLICT(follower_id, following_id)
		DO UPDATE SET status = ?
	`

	_, err := r.DB.Exec(
		query,
		followerID,
		followingID,
		status,
		status,
	)

	return err
}

// Unfollow removes the follow relationship.
// Works for both pending and accepted follows.
func (r *FollowRepository) Unfollow(followerID, followingID int) error {
	query := `
		DELETE FROM FOLLOWS
		WHERE follower_id = ?
		AND following_id = ?
	`

	_, err := r.DB.Exec(
		query,
		followerID,
		followingID,
	)

	return err
}

// AcceptFollow changes a pending request into an accepted follow.
// Returns an error if no pending follow request was found.
func (r *FollowRepository) AcceptFollow(followerID, followingID int) error {
	query := `
		UPDATE FOLLOWS
		SET status = 'accepted'
		WHERE follower_id = ?
		AND following_id = ?
		AND status = 'pending'
	`

	res, err := r.DB.Exec(
		query,
		followerID,
		followingID,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no pending follow request found")
	}

	return nil
}

// RejectFollow removes a pending follow request.
// No history is stored.
// Returns an error if no pending follow request was found.
func (r *FollowRepository) RejectFollow(followerID, followingID int) error {
	query := `
		DELETE FROM FOLLOWS
		WHERE follower_id = ?
		AND following_id = ?
		AND status = 'pending'
	`

	res, err := r.DB.Exec(
		query,
		followerID,
		followingID,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no pending follow request found")
	}

	return nil
}

// GetFollowers returns all accepted followers of a user.
func (r *FollowRepository) GetFollowers(userID int) ([]User, error) {
	query := `
		SELECT 
			u.id,
			u.firstname,
			u.lastname,
			u.nickname,
			u.avatar
		FROM USERS u
		JOIN FOLLOWS f 
			ON u.id = f.follower_id
		WHERE f.following_id = ?
		AND f.status = 'accepted'
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.scanUsers(rows)
}

// GetFollowing returns all users followed by a user.
func (r *FollowRepository) GetFollowing(userID int) ([]User, error) {
	query := `
		SELECT 
			u.id,
			u.firstname,
			u.lastname,
			u.nickname,
			u.avatar
		FROM USERS u
		JOIN FOLLOWS f 
			ON u.id = f.following_id
		WHERE f.follower_id = ?
		AND f.status = 'accepted'
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.scanUsers(rows)
}

// scanUsers converts database rows into User objects.
func (r *FollowRepository) scanUsers(rows *sql.Rows) ([]User, error) {
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
		); err != nil {
			return nil, err
		}

		u.Nickname = nickname.String
		u.Avatar = avatar.String

		users = append(users, u)
	}

	return users, nil
}

// GetFollowStatus returns the current relationship status.
//
// Returns:
// accepted -> already following
// pending  -> request waiting
// none     -> no relationship
func (r *FollowRepository) GetFollowStatus(user1ID, user2ID int) (string, error) {
	var status string

	query := `
		SELECT status
		FROM FOLLOWS
		WHERE follower_id = ?
		AND following_id = ?
	`

	err := r.DB.QueryRow(
		query,
		user1ID,
		user2ID,
	).Scan(&status)

	if err == sql.ErrNoRows {
		return "none", nil
	}

	return status, err
}
