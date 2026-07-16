package repository

import (
	"database/sql"
)

type Follow struct {
	FollowerID  int    `json:"follower_id"`
	FollowingID int    `json:"following_id"`
	Status      string `json:"status"`
}

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) RequestFollow(followerID, followingID int, status string) error {
	query := `INSERT INTO FOLLOWS (follower_id, following_id, status) VALUES (?, ?, ?)
	          ON CONFLICT(follower_id, following_id) DO UPDATE SET status = ?`
	_, err := r.DB.Exec(query, followerID, followingID, status, status)
	return err
}

func (r *FollowRepository) UpdateStatus(followerID, followingID int, status string) error {
	query := `UPDATE FOLLOWS SET status = ? WHERE follower_id = ? AND following_id = ?`
	_, err := r.DB.Exec(query, status, followerID, followingID)
	return err
}

func (r *FollowRepository) Unfollow(followerID, followingID int) error {
	query := `DELETE FROM FOLLOWS WHERE follower_id = ? AND following_id = ?`
	_, err := r.DB.Exec(query, followerID, followingID)
	return err
}

func (r *FollowRepository) GetFollowers(userID string) ([]User, error) {
	query := `SELECT u.id, u.firstname, u.lastname, u.nickname, u.avatar FROM USERS u
	          JOIN FOLLOWS f ON u.id = f.follower_id WHERE f.following_id = ? AND f.status = 'accepted'`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanUsers(rows)
}

func (r *FollowRepository) GetFollowing(userID string) ([]User, error) {
	query := `SELECT u.id, u.firstname, u.lastname, u.nickname, u.avatar FROM USERS u
	          JOIN FOLLOWS f ON u.id = f.following_id WHERE f.follower_id = ? AND f.status = 'accepted'`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanUsers(rows)
}

func (r *FollowRepository) scanUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	for rows.Next() {
		var u User
		var nickname, avatar sql.NullString
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &nickname, &avatar); err != nil {
			return nil, err
		}
		u.Nickname = nickname.String
		u.Avatar = avatar.String
		users = append(users, u)
	}
	return users, nil
}

// true if user1 follow user2
func (r *FollowRepository) IsFollowing(user1ID, user2ID int) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM FOLLOWS
			WHERE follower_id = ?
			AND following_id = ?
			AND status = 'accepted'
		)
	`

	err := r.DB.QueryRow(query, user1ID, user2ID).Scan(&exists)

	return exists, err
}
