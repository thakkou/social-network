package repository

import (
	"database/sql"
	"time"
)

// Post represents a post stored in the database
type Post struct {
	ID        int       `json:"id"`         // Unique post identifier
	UserID    int       `json:"user_id"`    // User who created the post
	CreatedAt time.Time `json:"created_at"` // Creation timestamp
	Title     string    `json:"title"`      // Post title
	Text      string    `json:"text"`       // Post content
	Image     string    `json:"image"`      // Optional image path
	Privacy   string    `json:"privacy"`    // public, almost_private, private
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `json:"id"`         // Unique comment identifier
	UserID    int       `json:"user_id"`    // User who created the comment
	PostID    int       `json:"post_id"`    // Post being commented on
	CreatedAt time.Time `json:"created_at"` // Comment creation timestamp
	Text      string    `json:"text"`       // Comment content
}

// PostRepository handles all database operations related to posts
type PostRepository struct {
	DB *sql.DB
}

// Creates a new PostRepository instance
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// CreatePost inserts a new post into the database.
// It also handles:
// - private post allowed users
// - post categories
// The operation uses a transaction to ensure all inserts succeed together.
func (r *PostRepository) CreatePost(p *Post, allowedUserIDs []int, categories []int) error {
	// Start database transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Rollback automatically if something fails
	defer tx.Rollback()

	// Insert the main post information
	query := `
		INSERT INTO POSTS 
		(user_id, created_at, title, text, image, privacy) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := tx.Exec(
		query,
		p.UserID,
		p.CreatedAt.Format("2006-01-02 15:04:05"),
		p.Title,
		p.Text,
		p.Image,
		p.Privacy,
	)
	if err != nil {
		return err
	}

	// Get the generated post ID
	postID, _ := res.LastInsertId()
	p.ID = int(postID)

	// If the post is private, store users allowed to view it
	if p.Privacy == "private" {
		for _, uid := range allowedUserIDs {

			_, err = tx.Exec(
				`INSERT INTO POST_ALLOWED_USERS (post_id, user_id) VALUES (?, ?)`,
				p.ID,
				uid,
			)
			if err != nil {
				return err
			}
		}
	}

	// Link the post with its categories
	for _, catID := range categories {

		_, err = tx.Exec(
			`INSERT INTO POST_CATEGORY (post_id, category_id) VALUES (?, ?)`,
			p.ID,
			catID,
		)
		if err != nil {
			return err
		}
	}

	// Save all changes permanently
	return tx.Commit()
}

// GetVisiblePosts returns all posts that the given user has permission to see.
// Visibility rules:
// - Public posts are visible to everyone
// - Users can see their own posts
// - Almost private posts are visible to accepted followers
// - Private posts are visible only to allowed users
func (r *PostRepository) GetVisiblePosts(userID int) ([]Post, error) {
	query := `
		SELECT DISTINCT 
			p.id,
			p.user_id,
			p.created_at,
			p.title,
			p.text,
			p.image,
			p.privacy
		FROM POSTS p

		LEFT JOIN FOLLOWS f 
			ON p.user_id = f.following_id 
			AND f.follower_id = ?
			AND f.status = 'accepted'

		LEFT JOIN POST_ALLOWED_USERS pau 
			ON p.id = pau.post_id

		WHERE p.privacy = 'public'
		   OR p.user_id = ?
		   OR (p.privacy = 'almost_private' AND f.status = 'accepted')
		   OR (p.privacy = 'private' AND pau.user_id = ?)

		ORDER BY p.created_at DESC
	`

	rows, err := r.DB.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	// Convert database rows into Post objects
	for rows.Next() {

		var p Post
		var title, text, image sql.NullString
		var tStr string

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&tStr,
			&title,
			&text,
			&image,
			&p.Privacy,
		)
		if err != nil {
			return nil, err
		}

		// Convert stored string date back to time.Time
		p.CreatedAt, _ = time.Parse(
			"2006-01-02 15:04:05",
			tStr,
		)

		// Handle nullable database fields
		p.Title = title.String
		p.Text = text.String
		p.Image = image.String

		posts = append(posts, p)
	}

	return posts, nil
}

// AddComment creates a new comment for a post.
func (r *PostRepository) AddComment(c *Comment) error {
	query := `
		INSERT INTO COMMENTS 
		(user_id, post_id, created_at, text)
		VALUES (?, ?, ?, ?)
	`

	res, err := r.DB.Exec(
		query,
		c.UserID,
		c.PostID,
		c.CreatedAt.Format("2006-01-02 15:04:05"),
		c.Text,
	)
	if err != nil {
		return err
	}

	// Store generated comment ID
	id, _ := res.LastInsertId()
	c.ID = int(id)

	return nil
}

// SetPostReaction creates or updates a user's reaction on a post.
// isLike values:
//
//	1 -> like
//	0 -> remove/neutral
//
// -1 -> dislike
func (r *PostRepository) SetPostReaction(userID, postID, isLike int) error {
	query := `
		INSERT INTO POST_REACTIONS 
		(user_id, post_id, is_like)
		VALUES (?, ?, ?)

		ON CONFLICT(user_id, post_id)
		DO UPDATE SET is_like = ?
	`

	_, err := r.DB.Exec(
		query,
		userID,
		postID,
		isLike,
		isLike,
	)

	return err
}

func (r *PostRepository) GetPostsUserID(userID int) ([]Post, error) {
	// Select all posts created by the given user
	query := `
		SELECT 
			id,
			user_id,
			created_at,
			title,
			text,
			image,
			privacy
		FROM POSTS
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	// Convert database rows into Post objects
	for rows.Next() {

		var p Post
		var title, text, image sql.NullString
		var createdAt string

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&createdAt,
			&title,
			&text,
			&image,
			&p.Privacy,
		)
		if err != nil {
			return nil, err
		}

		// Convert stored date string into time.Time
		p.CreatedAt, _ = time.Parse(
			"2006-01-02 15:04:05",
			createdAt,
		)

		// Handle nullable fields
		p.Title = title.String
		p.Text = text.String
		p.Image = image.String

		posts = append(posts, p)
	}

	// Check if iteration caused an error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
