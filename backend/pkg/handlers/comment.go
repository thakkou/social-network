package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	dblayer "01social/pkg/models/db_layer"
	"01social/pkg/repository"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

// CreateComment adds a new comment to a post.
// @Summary Create a comment
// @Description Adds a comment to a post. Accepts JSON or multipart/form-data with optional image.
// @Tags Comments
// @Accept json
// @Produce json
// @Param postId body int false "Post ID (for JSON requests)"
// @Param text body string false "Comment text"
// @Param postId formData integer false "Post ID (for multipart requests)"
// @Param text formData string false "Comment body"
// @Param image formData file false "Optional comment image"
// @Success 201 {object} map[string]any "Comment created successfully"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/comments/create [post]
func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/comments/create" {
		utilities.WriteJSON(w, http.StatusNotFound, "page not found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	type CommentReq struct {
		PostId any    `json:"postId"`
		Text   string `json:"text"`
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	var text string
	var postID int
	var imagePath string

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid form data", nil)
			return
		}

		pidStr := r.FormValue("postId")
		pid, err := strconv.Atoi(pidStr)
		if err != nil || pid <= 0 {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
			return
		}
		postID = pid

		text = strings.TrimSpace(r.FormValue("text"))

		if file, header, err := r.FormFile("image"); err == nil {
			defer file.Close()
			if saved, saveErr := utilities.SaveImage(file, header, "uploads/comments/"); saveErr == nil {
				imagePath = saved
			}
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		commentReq, err := utilities.ReadJSONRequest[CommentReq](r)
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}

		pid, err := utilities.ToInt(commentReq.PostId)
		if err != nil || pid <= 0 {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
			return
		}
		postID = pid

		text = strings.TrimSpace(commentReq.Text)
	} else {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json or multipart/form-data", nil)
		return
	}

	if text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "comment cannot be empty", nil)
		return
	}
	if len(text) > 1000 {
		utilities.WriteJSON(w, http.StatusBadRequest, "comment cannot exceed 1000 characters", nil)
		return
	}

	user, err := Repos.User.GetByID(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "user not found", nil)
		return
	}

	comment := &repository.Comment{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now().UTC(),
		Text:      text,
		Image:     imagePath,
	}

	if err := Repos.Comment.AddComment(comment); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not create comment", nil)
		return
	}

	type Res struct {
		ID        int       `json:"id"`
		Text      string    `json:"text"`
		PostID    int       `json:"postId"`
		UserID    int       `json:"userId"`
		CreatedAt time.Time `json:"createdAt"`
		Nickname  string    `json:"nickname"`
		Image     string    `json:"image"`
	}

	res := Res{
		ID:        comment.ID,
		Text:      comment.Text,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		CreatedAt: comment.CreatedAt,
		Nickname:  user.Nickname,
		Image:     comment.Image,
	}

	// Notify the post author about the new comment via WS
	postAuthor, err := Repos.Post.GetPostAuthor(postID)
	if err == nil && postAuthor != userID {
		ws.NotifyUser(strconv.Itoa(postAuthor), "new_comments", map[string]any{
			"post_id":    postID,
			"comment_id": comment.ID,
			"user_id":    userID,
			"nickname":   user.Nickname,
			"text":       comment.Text[:min(len(comment.Text), 80)],
		})

		Repos.Notification.Create(&repository.Notification{
			UserID:     postAuthor,
			ActorID:    userID,
			Type:       "comment",
			ObjectType: "post",
			ObjectID:   postID,
		})
	}

	utilities.WriteJSON(w, http.StatusCreated, "comment created successfully", res)
}

// CommentResolver handles comment reactions and deletion.
// @Summary Comment operations
// @Description Handles like/dislike on comments and comment deletion.
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Router /api/comments/{id}/like [post]
// @Router /api/comments/{id}/dislike [post]
// @Router /api/comments/{id}/delete [delete]
func CommentResolver(w http.ResponseWriter, r *http.Request) {
	segments := getPathSegments(r)
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "comments" {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	commentID, err := strconv.Atoi(segments[2])
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid comment id", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	endpoint := ""
	if len(segments) >= 4 {
		endpoint = segments[3]
	}

	switch endpoint {
	case "like", "dislike":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		isLike := 1
		if endpoint == "dislike" {
			isLike = -1
		}

		if err := Repos.Reaction.SetCommentReaction(userID, commentID, isLike); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not react to comment", nil)
			return
		}

		reactionCounts, err := Repos.Reaction.GetReactionComment(commentID, userID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not get reactions", nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, endpoint+"d", map[string]any{
			"comment_id": commentID,
			"likes":      reactionCounts.Likes,
			"dislikes":   reactionCounts.Dislikes,
			"is_liked":   reactionCounts.IsLiked,
		})

	case "delete":
		if r.Method != http.MethodDelete {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if err := Repos.Comment.DeleteComment(commentID, userID); err != nil {
			utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "comment deleted successfully", nil)
	default:
		utilities.WriteJSON(w, http.StatusNotFound, "unknown endpoint", nil)
	}
}

// GetCommentsByPost
func GetCommentsByPost(postId int) ([]dblayer.Comment, error) {
	return GetCommentsByPostWithPagination(postId, 0, 0)
}

func GetCommentsByPostWithPagination(postId, limit, lastID int) ([]dblayer.Comment, error) {
	var comments []dblayer.Comment

	if limit <= 0 {
		limit = 10
	}

	query := `SELECT id, user_id, created_at, text, COALESCE(image, '') FROM Comments WHERE post_id = ?`
	args := []any{postId}

	if lastID > 0 {
		query += " AND id < ?"
		args = append(args, lastID)
	}

	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("getCommentsByPost error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c dblayer.Comment
		var createdAtStr string
		if err := rows.Scan(&c.Id, &c.UserId, &createdAtStr, &c.Text, &c.Image); err != nil {
			return nil, fmt.Errorf("getCommentsByPost scan error: %v", err)
		}

		c.Created_at, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

		// get username
		if err := db.Database.QueryRow(
			"SELECT u.nickname FROM users u INNER JOIN comments c ON c.user_id = u.id WHERE c.id = ?",
			c.Id,
		).Scan(&c.Nickname); err != nil {
			return nil, fmt.Errorf("getCommentsByPost username error: %v", err)
		}

		// get timeago
		c.TimeAgo = utilities.TimeAgo(c.Created_at)

		// get reactions
		if c.LikeCount, c.DislikeCount, err = GetReactionsByComment(c.Id); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCommentsByPost rows error: %v", err)
	}
	return comments, nil
}

// DeleteComment
func DeleteComment(commentId, userId int) error {
	tx, err := db.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var dbUserId int
	err = tx.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentId).Scan(&dbUserId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return err
	}
	if dbUserId != userId {
		return fmt.Errorf("not your comment")
	}

	_, err = tx.Exec("DELETE FROM comments WHERE id = ?", commentId)
	if err != nil {
		return err
	}
	return tx.Commit()
}
