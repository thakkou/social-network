package repository

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	Privacy   string    `json:"privacy"`
}

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) CreatePost(p *Post, allowedUserIDs []int, categories []int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO POSTS (user_id, created_at, title, text, image, privacy) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(query, p.UserID, p.CreatedAt.Format("2006-01-02 15:04:05"), p.Title, p.Text, p.Image, p.Privacy)
	if err != nil {
		return err
	}
	postID, _ := res.LastInsertId()
	p.ID = int(postID)

	if p.Privacy == "private" {
		for _, uid := range allowedUserIDs {
			_, err = tx.Exec(`INSERT INTO POST_ALLOWED_USERS (post_id, user_id) VALUES (?, ?)`, p.ID, uid)
			if err != nil {
				return err
			}
		}
	}

	for _, catID := range categories {
		_, err = tx.Exec(`INSERT INTO POST_CATEGORY (post_id, category_id) VALUES (?, ?)`, p.ID, catID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostRepository) GetVisiblePosts(userID int) ([]Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text, p.image, p.privacy FROM POSTS p
		LEFT JOIN FOLLOWS f ON p.user_id = f.following_id AND f.follower_id = ? AND f.status = 'accepted'
		LEFT JOIN POST_ALLOWED_USERS pau ON p.id = pau.post_id
		WHERE p.privacy = 'public' 
		   OR p.user_id = ?
		   OR (p.privacy = 'almost_private' AND f.status = 'accepted')
		   OR (p.privacy = 'private' AND pau.user_id = ?)
		ORDER BY p.created_at DESC`

	rows, err := r.DB.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var title, text, image sql.NullString
		var tStr string
		if err := rows.Scan(&p.ID, &p.UserID, &tStr, &title, &text, &image, &p.Privacy); err != nil {
			return nil, err
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", tStr)
		p.Title = title.String
		p.Text = text.String
		p.Image = image.String
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) AddComment(c *Comment) error {
	query := `INSERT INTO COMMENTS (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, c.UserID, c.PostID, c.CreatedAt.Format("2006-01-02 15:04:05"), c.Text)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	c.ID = int(id)
	return nil
}

func (r *PostRepository) SetPostReaction(userID, postID, isLike int) error {
	query := `INSERT INTO POST_REACTIONS (user_id, post_id, is_like) VALUES (?, ?, ?)
	          ON CONFLICT(user_id, post_id) DO UPDATE SET is_like = ?`
	_, err := r.DB.Exec(query, userID, postID, isLike, isLike)
	return err
}
