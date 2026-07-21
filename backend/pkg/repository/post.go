package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Post represents a post stored in the database
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	Privacy   string    `json:"privacy"`
	// New fields to hold your metadata
	LikeCount    int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
	IsLiked      int `json:"is_liked"` // 1 = liked, -1 = disliked, 0 = neutral
	CommentCount int `json:"comment_count"`
}

// PostRepository handles all database operations related to posts
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new PostRepository instance
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// CreatePost inserts a new post into the database.
// It also handles:
// - private post allowed users
// - post categories
// The operation uses a transaction to ensure all inserts succeed together.
func (r *PostRepository) CreatePost(p *Post, allowedUserIDs []int, categoryIDs []int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	postID, _ := res.LastInsertId()
	p.ID = int(postID)

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

	for _, catID := range categoryIDs {
		_, err = tx.Exec(
			`INSERT INTO POST_CATEGORY (post_id, category_id) VALUES (?, ?)`,
			p.ID,
			catID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// EnrichPostMetadata populates counts and current user's reaction status
func (r *PostRepository) EnrichPostMetadata(userID int, p *Post) error {
	likes, dislikes, err := r.GetReactionCounts(p.ID)
	if err != nil {
		return err
	}
	p.LikeCount = likes
	p.DislikeCount = dislikes

	comments, err := r.GetCommentCount(p.ID)
	if err != nil {
		return err
	}
	p.CommentCount = comments

	if userID > 0 {
		reaction, err := r.GetUserReaction(userID, p.ID)
		if err != nil {
			return err
		}
		p.IsLiked = reaction
	}

	return nil
}

// scanPostRow reads a single POSTS row (id, user_id, created_at, title, text, image, privacy)
func scanPostRow(scanner interface{ Scan(...any) error }) (Post, error) {
	var p Post
	var title, text, image sql.NullString
	var createdAt string

	err := scanner.Scan(&p.ID, &p.UserID, &createdAt, &title, &text, &image, &p.Privacy)
	if err != nil {
		return p, err
	}

	// Try multiple formats if your DB formats vary, but make sure to capture the result
	parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		p.CreatedAt = parsedTime
	} else {
		// Fallback or handle standard RFC3339 if needed
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	}

	p.Title = title.String
	p.Text = text.String
	p.Image = image.String

	return p, nil
}

// GetVisiblePosts returns all posts that the given user has permission to see.
func (r *PostRepository) GetVisiblePosts(userID int) ([]Post, error) {
	query := `
		SELECT DISTINCT 
			p.id, p.user_id, p.created_at, p.title, p.text, p.image, p.privacy
		FROM POSTS p
		LEFT JOIN FOLLOWS f 
			ON p.user_id = f.following_id 
			AND f.follower_id = ?
			AND f.status = 'accepted'
		LEFT JOIN POST_ALLOWED_USERS pau 
			ON p.id = pau.post_id
			AND pau.user_id = ?
		WHERE p.privacy = 'public'
		   OR p.user_id = ?
		   OR (p.privacy = 'almost_private' AND f.status = 'accepted')
		   OR (p.privacy = 'private' AND pau.user_id = ?)
		ORDER BY p.created_at DESC
	`

	rows, err := r.DB.Query(query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		p, err := scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (r *PostRepository) GetVisiblePostByID(userID, postID int) (*Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text, p.image, p.privacy
		FROM POSTS p
		LEFT JOIN FOLLOWS f 
			ON p.user_id = f.following_id 
			AND f.follower_id = ?
			AND f.status = 'accepted'
		LEFT JOIN POST_ALLOWED_USERS pau 
			ON p.id = pau.post_id
			AND pau.user_id = ?
		WHERE p.id = ?
		  AND (
			p.privacy = 'public'
			OR p.user_id = ?
			OR (p.privacy = 'almost_private' AND f.status = 'accepted')
			OR (p.privacy = 'private' AND pau.user_id = ?)
		)
	`

	row := r.DB.QueryRow(query, userID, userID, postID, userID, userID)
	p, err := scanPostRow(row)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPostByID fetches a single post by its ID.
func (r *PostRepository) GetPostByID(id int) (*Post, error) {
	row := r.DB.QueryRow(`
		SELECT id, user_id, created_at, title, text, image, privacy
		FROM POSTS
		WHERE id = ?
	`, id)

	p, err := scanPostRow(row)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPostsUserID returns all posts created by the given user.
func (r *PostRepository) GetPostsUserID(userID int) ([]Post, error) {
	query := `
		SELECT id, user_id, created_at, title, text, image, privacy
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
	for rows.Next() {
		p, err := scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		fmt.Println("post befaure", p.LikeCount, p.DislikeCount, p.IsLiked)
		// Enrich the post with metadata stats
		if err := r.EnrichPostMetadata(userID, &p); err != nil {
			fmt.Println("error enrishing post")
			return nil, err
		}
		fmt.Println("post after", p.LikeCount, p.DislikeCount, p.IsLiked)

		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// GetFilteredPosts returns posts filtered by category names, "liked by me" and
// "posted by me", paginated via a keyset (lastID) cursor.
func (r *PostRepository) GetFilteredPosts(
	userID int,
	categories []string,
	likedByMe, postedByMe bool,
	limit, lastID int,
) ([]Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text, p.image, p.privacy
		FROM POSTS p
		LEFT JOIN FOLLOWS f
			ON p.user_id = f.following_id
			AND f.follower_id = ?
			AND f.status = 'accepted'
		LEFT JOIN POST_ALLOWED_USERS pau
			ON p.id = pau.post_id
			AND pau.user_id = ?
		LEFT JOIN POST_CATEGORY pc ON p.id = pc.post_id
		LEFT JOIN CATEGORY c ON pc.category_id = c.id
	`

	var cond []string
	var args []any
	args = append(args, userID, userID)

	visibilityCond := `(
		p.privacy = 'public'
		OR p.user_id = ?
		OR (p.privacy = 'almost_private' AND f.status = 'accepted')
		OR (p.privacy = 'private' AND pau.user_id = ?)
	)`
	cond = append(cond, visibilityCond)
	args = append(args, userID, userID)

	if len(categories) > 0 {
		ph := make([]string, len(categories))
		for i, c := range categories {
			ph[i] = "?"
			args = append(args, c)
		}
		cond = append(cond, "c.name IN ("+strings.Join(ph, ",")+")")
	}

	if postedByMe {
		cond = append(cond, "p.user_id = ?")
		args = append(args, userID)
	}

	if likedByMe {
		query += " JOIN POST_REACTIONS pr ON p.id = pr.post_id "
		cond = append(cond, "pr.user_id = ? AND pr.is_like = 1")
		args = append(args, userID)
	}

	if lastID > 0 {
		cond = append(cond, "p.id < ?")
		args = append(args, lastID)
	}

	if len(cond) > 0 {
		query += " WHERE " + strings.Join(cond, " AND ")
	}

	query += " ORDER BY p.created_at DESC, p.id DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		p, err := scanPostRow(rows)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// DeletePost removes a post if userID is its owner.
func (r *PostRepository) DeletePost(postID, userID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var owner int
	err = tx.QueryRow("SELECT user_id FROM POSTS WHERE id = ?", postID).Scan(&owner)
	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}
	if err != nil {
		return err
	}
	if owner != userID {
		return fmt.Errorf("not your post")
	}

	if _, err := tx.Exec("DELETE FROM POSTS WHERE id = ?", postID); err != nil {
		return err
	}

	return tx.Commit()
}

// SetPostReaction creates or updates a user's reaction on a post.
// isLike values: 1 -> like, 0 -> remove/neutral, -1 -> dislike
func (r *PostRepository) SetPostReaction(userID, postID, isLike int) error {
	query := `
		INSERT INTO POST_REACTIONS 
		(user_id, post_id, is_like)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id)
		DO UPDATE SET is_like = ?
	`

	_, err := r.DB.Exec(query, userID, postID, isLike, isLike)
	return err
}

// GetReactionCounts returns the like/dislike counts for a post.
func (r *PostRepository) GetReactionCounts(postID int) (likes, dislikes int, err error) {
	err = r.DB.QueryRow(`
		SELECT
			COUNT(CASE WHEN is_like = 1 THEN 1 END),
			COUNT(CASE WHEN is_like = -1 THEN 1 END)
		FROM POST_REACTIONS
		WHERE post_id = ?
	`, postID).Scan(&likes, &dislikes)
	return likes, dislikes, err
}

// GetUserReaction returns the given user's reaction to a post (1, 0, or -1).
func (r *PostRepository) GetUserReaction(userID, postID int) (int, error) {
	var reaction int
	err := r.DB.QueryRow(`
		SELECT is_like FROM POST_REACTIONS WHERE user_id = ? AND post_id = ?
	`, userID, postID).Scan(&reaction)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	return reaction, err
}

func (r *PostRepository) GetCommentCount(postID int) (int, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM COMMENTS WHERE post_id = ?`, postID).Scan(&count)
	return count, err
}

// GetPostAuthor returns the user_id of the post author (the creator of the post).
func (r *PostRepository) GetPostAuthor(postID int) (int, error) {
	var authorID int
	err := r.DB.QueryRow(`SELECT user_id FROM POSTS WHERE id = ?`, postID).Scan(&authorID)
	if err != nil {
		return 0, err
	}
	return authorID, nil
}
