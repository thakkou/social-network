package repository

import (
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

// PostEngagement holds the enrichment data attached to a post in the feed:
// like/dislike counts, whether the requesting user liked it, and comment count.
type PostEngagement struct {
	LikesCount    int  `json:"likes_count"`
	DislikesCount int  `json:"dislikes_count"`
	IsLiked       bool `json:"is_liked"`
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

	// Whether the requesting user liked each post.
	likedQuery := `
SELECT group_post_id
FROM GROUP_POST_REACTIONS
WHERE user_id = ? AND is_like = 1 AND group_post_id IN (` + placeholders + `)
`
	likedArgs := append([]interface{}{userID}, args...)
	likedRows, err := r.DB.Query(likedQuery, likedArgs...)
	if err != nil {
		return nil, err
	}
	defer likedRows.Close()

	for likedRows.Next() {
		var postID int
		if err := likedRows.Scan(&postID); err != nil {
			return nil, err
		}
		e := engagement[postID]
		e.IsLiked = true
		engagement[postID] = e
	}
	if err := likedRows.Err(); err != nil {
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
