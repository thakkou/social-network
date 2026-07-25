package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
}

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// AddComment inserts a new comment for a post.
func (r *CommentRepository) AddComment(c *Comment) error {
	res, err := r.DB.Exec(
		`INSERT INTO COMMENTS (user_id, post_id, created_at, text, image) VALUES (?, ?, ?, ?, ?)`,
		c.UserID,
		c.PostID,
		c.CreatedAt.Format("2006-01-02 15:04:05"),
		c.Text,
		sql.NullString{String: c.Image, Valid: c.Image != ""},
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		c.ID = int(id)
	}
	return err
}

// DeleteComment removes a comment if the user is the owner.
func (r *CommentRepository) DeleteComment(commentID, userID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var ownerID int
	err = tx.QueryRow("SELECT user_id FROM COMMENTS WHERE id = ?", commentID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return err
	}
	if ownerID != userID {
		return fmt.Errorf("not your comment")
	}

	if _, err := tx.Exec("DELETE FROM COMMENTS WHERE id = ?", commentID); err != nil {
		return err
	}

	return tx.Commit()
}

// GetCommentsByPost returns all comments for a post, newest first.
func (r *CommentRepository) GetCommentsByPost(postID int) ([]Comment, error) {
	return r.getComments(postID, 0, 0)
}

// GetCommentsByPostPaginated returns comments for a post using a keyset (lastID) cursor.
func (r *CommentRepository) GetCommentsByPostPaginated(postID, limit, lastID int) ([]Comment, error) {
	return r.getComments(postID, limit, lastID)
}

func (r *CommentRepository) getComments(postID, limit, lastID int) ([]Comment, error) {
	query := `SELECT id, user_id, post_id, created_at, text, COALESCE(image, '') FROM COMMENTS WHERE post_id = ?`
	args := []any{postID}

	if lastID > 0 {
		query += " AND id < ?"
		args = append(args, lastID)
	}

	query += " ORDER BY created_at DESC, id DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var createdAt string
		if err := rows.Scan(&c.ID, &c.UserID, &c.PostID, &createdAt, &c.Text, &c.Image); err != nil {
			return nil, err
		}

		// Try parsing the timestamp with the standard format first.
		// If that fails, try alternative formats (e.g. RFC3339 from JSON round-trips).
		var parsedTime time.Time
		var parseErr error
		parsedTime, parseErr = time.Parse("2006-01-02 15:04:05", createdAt)
		if parseErr != nil {
			parsedTime, parseErr = time.Parse(time.RFC3339, createdAt)
		}
		if parseErr != nil {
			parsedTime, parseErr = time.Parse("2006-01-02 15:04:05.000", createdAt)
		}
		if parseErr != nil {
			// Last resort: log and use current time (better than zero date)
			log.Printf("[COMMENTS] failed to parse created_at %q for comment %d, using current time", createdAt, c.ID)
			parsedTime = time.Now()
		}

		c.CreatedAt = parsedTime
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

// GetCommentCount returns the number of comments attached to a post.

// GetCommentByID returns a single comment by its ID.
func (r *CommentRepository) GetCommentByID(commentID int) (*Comment, error) {
	var c Comment
	var createdAt string

	err := r.DB.QueryRow(
		`SELECT id, user_id, post_id, created_at, text, COALESCE(image, '') 
		 FROM COMMENTS 
		 WHERE id = ?`,
		commentID,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.PostID,
		&createdAt,
		&c.Text,
		&c.Image,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}

	if err != nil {
		return nil, err
	}

	c.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
