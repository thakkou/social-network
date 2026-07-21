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

// SetPostReaction toggles a user's reaction on a post.
// - If no reaction → inserts it
// - If same reaction exists → deletes it (removes the reaction)
// - If different reaction exists → updates it
func (r *ReactionRepository) SetPostReaction(userID, postID, isLike int) error {
	var existing int
	err := r.DB.QueryRow(
		"SELECT is_like FROM POST_REACTIONS WHERE user_id = ? AND post_id = ?",
		userID, postID,
	).Scan(&existing)

	if err == sql.ErrNoRows {
		// No reaction → insert
		_, err = r.DB.Exec(
			"INSERT INTO POST_REACTIONS (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userID, postID, isLike,
		)
		return err
	}
	if err != nil {
		return err
	}

	if existing == isLike {
		// Same reaction → remove (delete)
		_, err = r.DB.Exec(
			"DELETE FROM POST_REACTIONS WHERE user_id = ? AND post_id = ?",
			userID, postID,
		)
		return err
	}

	// Different reaction → update
	_, err = r.DB.Exec(
		"UPDATE POST_REACTIONS SET is_like = ? WHERE user_id = ? AND post_id = ?",
		isLike, userID, postID,
	)
	return err
}

// SetCommentReaction toggles a user's reaction on a comment.
// - If no reaction → inserts it
// - If same reaction exists → deletes it (removes the reaction)
// - If different reaction exists → updates it
func (r *ReactionRepository) SetCommentReaction(userID, commentID, isLike int) error {
	var existing int
	err := r.DB.QueryRow(
		"SELECT is_like FROM COMMENT_REACTIONS WHERE user_id = ? AND comment_id = ?",
		userID, commentID,
	).Scan(&existing)

	if err == sql.ErrNoRows {
		// No reaction → insert
		_, err = r.DB.Exec(
			"INSERT INTO COMMENT_REACTIONS (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			userID, commentID, isLike,
		)
		return err
	}
	if err != nil {
		return err
	}

	if existing == isLike {
		// Same reaction → remove (delete)
		_, err = r.DB.Exec(
			"DELETE FROM COMMENT_REACTIONS WHERE user_id = ? AND comment_id = ?",
			userID, commentID,
		)
		return err
	}

	// Different reaction → update
	_, err = r.DB.Exec(
		"UPDATE COMMENT_REACTIONS SET is_like = ? WHERE user_id = ? AND comment_id = ?",
		isLike, userID, commentID,
	)
	return err
}
