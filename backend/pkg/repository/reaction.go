package repository

import (
	"database/sql"
)

type ReactionCounts struct {
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
	IsLiked  int `json:"is_liked"` // 1 = liked, -1 = disliked, 0 = no reaction
}
type ReactionRepository struct {
	DB *sql.DB
}

// NewReactionRepository creates a new ReactionRepository instance
func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{DB: db}
}

// GetReactionPost fetches aggregated counts AND the specific user's reaction state for a post
func (r *ReactionRepository) GetReactionPost(postID, userID int) (ReactionCounts, error) {
	var counts ReactionCounts

	// 1. Get Totals
	err := r.DB.QueryRow(`
        SELECT
            COUNT(CASE WHEN is_like = 1 THEN 1 END),
            COUNT(CASE WHEN is_like = -1 THEN 1 END)
        FROM POST_REACTIONS
        WHERE post_id = ?
    `, postID).Scan(&counts.Likes, &counts.Dislikes)
	if err != nil {
		return counts, err
	}

	// 2. Find if this specific user reacted to it
	err = r.DB.QueryRow(`
        SELECT is_like FROM POST_REACTIONS WHERE user_id = ? AND post_id = ?
    `, userID, postID).Scan(&counts.IsLiked)

	if err == sql.ErrNoRows {
		counts.IsLiked = 0 // User hasn't reacted yet
		err = nil
	}

	return counts, err
}

// GetReactionComment fetches aggregated counts AND the specific user's reaction state for a comment
func (r *ReactionRepository) GetReactionComment(commentID, userID int) (ReactionCounts, error) {
	var counts ReactionCounts

	// 1. Get Totals
	err := r.DB.QueryRow(`
        SELECT
            COUNT(CASE WHEN is_like = 1 THEN 1 END),
            COUNT(CASE WHEN is_like = -1 THEN 1 END)
        FROM COMMENT_REACTIONS
        WHERE comment_id = ?
    `, commentID).Scan(&counts.Likes, &counts.Dislikes)
	if err != nil {
		return counts, err
	}

	// 2. Find if this specific user reacted to it
	err = r.DB.QueryRow(`
        SELECT is_like FROM COMMENT_REACTIONS WHERE user_id = ? AND comment_id = ?
    `, userID, commentID).Scan(&counts.IsLiked)

	if err == sql.ErrNoRows {
		counts.IsLiked = 0 // User hasn't reacted yet
		err = nil
	}

	return counts, err
}

// SetPostReaction creates, updates, or toggles a user's reaction on a post
// isLike values: 1 -> like, 0 -> clear/neutral, -1 -> dislike
func (r *ReactionRepository) SetPostReaction(userID, postID, isLike int) error {
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

// SetCommentReaction creates, updates, or toggles a user's reaction on a comment
// isLike values: 1 -> like, 0 -> clear/neutral, -1 -> dislike
func (r *ReactionRepository) SetCommentReaction(userID, commentID, isLike int) error {
	query := `
        INSERT INTO COMMENT_REACTIONS 
        (user_id, comment_id, is_like)
        VALUES (?, ?, ?)
        ON CONFLICT(user_id, comment_id)
        DO UPDATE SET is_like = ?
    `
	_, err := r.DB.Exec(query, userID, commentID, isLike, isLike)
	return err
}
