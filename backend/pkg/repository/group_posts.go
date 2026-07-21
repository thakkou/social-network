package repository

import (
	"database/sql"
	"fmt"
	"time"

	"01social/pkg/utilities"
)

type GroupPost struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
}

type GroupPostComment struct {
	ID          int       `json:"id"`
	GroupPostID int       `json:"group_post_id"`
	UserID      int       `json:"user_id"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupPostCommentWithAuthor struct {
	ID          int    `json:"id"`
	GroupPostID int    `json:"group_post_id"`
	UserID      int    `json:"user_id"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at"`
	Nickname    string `json:"nickname"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Avatar      string `json:"avatar"`
	LikesCount  int    `json:"likes_count"`
	DislikesCount int  `json:"dislikes_count"`
	IsLiked     int    `json:"is_liked"` // 1 liked, -1 disliked, 0 none
}

func (r *GroupRepository) ListGroupPosts(groupID int) ([]GroupPost, error) {
	rows, err := r.DB.Query(`SELECT id, group_id, user_id, created_at, title, text, image FROM GROUP_POSTS WHERE group_id = ? ORDER BY created_at DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]GroupPost, 0)
	for rows.Next() {
		var p GroupPost
		var createdAt string
		if err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &createdAt, &p.Title, &p.Text, &p.Image); err != nil {
			return nil, err
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *GroupRepository) ListGroupPostComments(groupPostID int) ([]GroupPostComment, error) {
	rows, err := r.DB.Query(`SELECT id, group_post_id, user_id, text, created_at FROM GROUP_POST_COMMENTS WHERE group_post_id = ? ORDER BY created_at ASC`, groupPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]GroupPostComment, 0)
	for rows.Next() {
		var c GroupPostComment
		var createdAt string
		if err := rows.Scan(&c.ID, &c.GroupPostID, &c.UserID, &c.Text, &createdAt); err != nil {
			return nil, err
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// GetGroupPostCommentsWithAuthors batch-fetches comments with user profile info for a set of post IDs.
// Returns a map keyed by post ID.
func (r *GroupRepository) GetGroupPostCommentsWithAuthors(postIDs []int, userID int) (map[int][]GroupPostCommentWithAuthor, error) {
	result := make(map[int][]GroupPostCommentWithAuthor, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}

	placeholders, args := utilities.PlaceholdersForInts(postIDs)

	query := `
SELECT
    gpc.id,
    gpc.group_post_id,
    gpc.user_id,
    gpc.text,
    gpc.created_at,
    COALESCE(u.nickname, '') AS nickname,
    u.firstname,
    u.lastname,
    COALESCE(u.avatar, '') AS avatar
FROM GROUP_POST_COMMENTS gpc
JOIN USERS u ON u.id = gpc.user_id
WHERE gpc.group_post_id IN (` + placeholders + `)
ORDER BY gpc.created_at ASC
`
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect comment IDs for batch reaction fetch.
	commentIDs := make([]int, 0)
	comments := make([]GroupPostCommentWithAuthor, 0)
	for rows.Next() {
		var c GroupPostCommentWithAuthor
		var postID int
		if err := rows.Scan(&c.ID, &postID, &c.UserID, &c.Text, &c.CreatedAt, &c.Nickname, &c.Firstname, &c.Lastname, &c.Avatar); err != nil {
			return nil, err
		}
		// Store by postID and collect commentID
		result[postID] = append(result[postID], c)
		commentIDs = append(commentIDs, c.ID)
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Batch-fetch reactions for all comments (totals)
	if len(commentIDs) > 0 {
		reactionPlaceholders, reactionArgs := utilities.PlaceholdersForInts(commentIDs)

		likesQuery := `
SELECT group_post_comment_id, SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END), SUM(CASE WHEN is_like = -1 THEN 1 ELSE 0 END)
FROM GROUP_POST_COMMENT_REACTIONS
WHERE group_post_comment_id IN (` + reactionPlaceholders + `)
GROUP BY group_post_comment_id
`
		likesRows, err := r.DB.Query(likesQuery, reactionArgs...)
		if err == nil {
			defer likesRows.Close()
			for likesRows.Next() {
				var cID, likes, dislikes int
				if err := likesRows.Scan(&cID, &likes, &dislikes); err != nil {
					continue
				}
				for postID := range result {
					for i := range result[postID] {
						if result[postID][i].ID == cID {
							result[postID][i].LikesCount = likes
							result[postID][i].DislikesCount = dislikes
						}
					}
				}
			}
		}

		// Get user's own reactions
		userReactionQuery := `
SELECT group_post_comment_id, is_like
FROM GROUP_POST_COMMENT_REACTIONS
WHERE user_id = ? AND group_post_comment_id IN (` + reactionPlaceholders + `)
`
		userArgs := append([]interface{}{userID}, commentIDsToInterface(commentIDs)...)
		userRows, err := r.DB.Query(userReactionQuery, userArgs...)
		if err == nil {
			defer userRows.Close()
			for userRows.Next() {
				var cID, isLike int
				if err := userRows.Scan(&cID, &isLike); err != nil {
					continue
				}
				for postID := range result {
					for i := range result[postID] {
						if result[postID][i].ID == cID {
							result[postID][i].IsLiked = isLike
						}
					}
				}
			}
		}
	}

	return result, nil
}

func commentIDsToInterface(ids []int) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}

func (r *GroupRepository) CreateGroupPost(gp *GroupPost) error {
	query := `INSERT INTO GROUP_POSTS (group_id, user_id, title, text, image) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, gp.GroupID, gp.UserID, gp.Title, gp.Text, gp.Image)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	gp.ID = int(id)
	return nil
}

func (r *GroupRepository) CreateGroupPostComment(groupPostID, userID int, text string) error {
	_, err := r.DB.Exec(`INSERT INTO GROUP_POST_COMMENTS (group_post_id, user_id, text) VALUES (?, ?, ?)`, groupPostID, userID, text)
	return err
}

// ReactToGroupPostComment creates, updates, or removes a reaction on a group post comment.
// isLike: 1 = like, -1 = dislike
func (r *GroupRepository) ReactToGroupPostComment(commentID, userID, isLike int) error {
	if isLike != 1 && isLike != -1 {
		return fmt.Errorf("invalid reaction")
	}

	var oldReaction int
	err := r.DB.QueryRow(
		"SELECT is_like FROM GROUP_POST_COMMENT_REACTIONS WHERE group_post_comment_id = ? AND user_id = ?",
		commentID, userID,
	).Scan(&oldReaction)

	if err == nil {
		// Same reaction -> remove (toggle off)
		if oldReaction == isLike {
			_, err = r.DB.Exec(
				"DELETE FROM GROUP_POST_COMMENT_REACTIONS WHERE group_post_comment_id = ? AND user_id = ?",
				commentID, userID,
			)
			return err
		}
		// Different reaction -> update
		_, err = r.DB.Exec(
			"UPDATE GROUP_POST_COMMENT_REACTIONS SET is_like = ? WHERE group_post_comment_id = ? AND user_id = ?",
			isLike, commentID, userID,
		)
		return err
	}

	if err != sql.ErrNoRows {
		return err
	}

	// No reaction -> insert
	_, err = r.DB.Exec(
		"INSERT INTO GROUP_POST_COMMENT_REACTIONS (group_post_comment_id, user_id, is_like) VALUES (?, ?, ?)",
		commentID, userID, isLike,
	)
	return err
}

// DeleteGroupPost deletes a group post if the user is the post author OR the group creator.
func (r *GroupRepository) DeleteGroupPost(postID, userID, creatorID int) error {
	// Check post exists and get its author
	var authorID int
	err := r.DB.QueryRow(`SELECT user_id FROM GROUP_POSTS WHERE id = ?`, postID).Scan(&authorID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}
	if err != nil {
		return err
	}
	if userID != authorID && userID != creatorID {
		return fmt.Errorf("not authorized to delete this post")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete reactions on this post
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_REACTIONS WHERE group_post_id = ?`, postID)
	// Delete comments on this post
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_COMMENTS WHERE group_post_id = ?`, postID)
	// Delete the post itself
	_, err = tx.Exec(`DELETE FROM GROUP_POSTS WHERE id = ?`, postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteGroupPostComment deletes a group post comment if the user is the comment author OR the group creator.
func (r *GroupRepository) DeleteGroupPostComment(commentID, userID, creatorID int) error {
	var authorID int
	err := r.DB.QueryRow(`SELECT user_id FROM GROUP_POST_COMMENTS WHERE id = ?`, commentID).Scan(&authorID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return err
	}
	if userID != authorID && userID != creatorID {
		return fmt.Errorf("not authorized to delete this comment")
	}

	_, err = r.DB.Exec(`DELETE FROM GROUP_POST_COMMENTS WHERE id = ?`, commentID)
	return err
}

// PostEngagement holds the enrichment data attached to a post in the feed:
// like/dislike counts, whether the requesting user liked it, and comment count.
// IsLiked: 1 = liked, -1 = disliked, 0 = no reaction
type PostEngagement struct {
	LikesCount    int  `json:"likes_count"`
	DislikesCount int  `json:"dislikes_count"`
	IsLiked       int  `json:"is_liked"`
	CommentsCount int  `json:"comments_count"`
}

// GetPostsEngagement batch-fetches engagement data for a set of post IDs,
// scoped to userID for the is_liked flag. Returns a map keyed by post ID
// with every requested ID present (zero-valued if it has no engagement yet).
func (r *GroupRepository) GetPostsEngagement(postIDs []int, userID int) (map[int]PostEngagement, error) {
	engagement := make(map[int]PostEngagement, len(postIDs))
	if len(postIDs) == 0 {
		return engagement, nil
	}
	for _, id := range postIDs {
		engagement[id] = PostEngagement{}
	}

	placeholders, args := utilities.PlaceholdersForInts(postIDs)

	// Like / dislike counts per post.
	reactionQuery := `
SELECT group_post_id, is_like, COUNT(*)
FROM GROUP_POST_REACTIONS
WHERE group_post_id IN (` + placeholders + `)
GROUP BY group_post_id, is_like
`
	reactionRows, err := r.DB.Query(reactionQuery, args...)
	if err != nil {
		return nil, err
	}
	defer reactionRows.Close()

	for reactionRows.Next() {
		var postID int
		var isLike int
		var count int
		if err := reactionRows.Scan(&postID, &isLike, &count); err != nil {
			return nil, err
		}
		e := engagement[postID]
		if isLike == 1 {
			e.LikesCount = count
		} else if isLike == -1 {
			e.DislikesCount = count
		}
		engagement[postID] = e
	}
	if err := reactionRows.Err(); err != nil {
		return nil, err
	}

	// Whether the requesting user liked or disliked each post.
	userReactionQuery := `
SELECT group_post_id, is_like
FROM GROUP_POST_REACTIONS
WHERE user_id = ? AND group_post_id IN (` + placeholders + `)
`
	userReactionArgs := append([]interface{}{userID}, args...)
	userRows, err := r.DB.Query(userReactionQuery, userReactionArgs...)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	for userRows.Next() {
		var postID int
		var isLike int
		if err := userRows.Scan(&postID, &isLike); err != nil {
			return nil, err
		}
		e := engagement[postID]
		e.IsLiked = isLike
		engagement[postID] = e
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	// Comment counts per post.
	commentQuery := `
SELECT group_post_id, COUNT(*)
FROM GROUP_POST_COMMENTS
WHERE group_post_id IN (` + placeholders + `)
GROUP BY group_post_id
`
	commentRows, err := r.DB.Query(commentQuery, args...)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	for commentRows.Next() {
		var postID int
		var count int
		if err := commentRows.Scan(&postID, &count); err != nil {
			return nil, err
		}
		e := engagement[postID]
		e.CommentsCount = count
		engagement[postID] = e
	}

	return engagement, commentRows.Err()
}
