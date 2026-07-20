package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"01social/pkg/middlewares"
	dblayer "01social/pkg/models/db_layer"
	"01social/pkg/repository"
	"01social/pkg/utilities"
)

// =========================
// REPOSITORIES
// =========================
//
// Call InitRepositories() once at startup, after db.Database has been
// opened/assigned (e.g. right after sql.Open in main()). Building these at
// package-var init time would risk capturing a nil *sql.DB if db.Database
// isn't set until later.

func getPathSegments(r *http.Request) []string {
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// =========================
// CORE POST ENRICHMENT
// =========================

func enrichPost(p *dblayer.Post, userId int) error {
	p.TimeAgo = utilities.TimeAgo(p.Created_at)

	user, err := Repos.User.GetByID(p.UserId)
	if err != nil {
		return err
	}
	p.Nickname = user.Nickname

	p.LikeCount, p.DislikeCount, err = Repos.Post.GetReactionCounts(p.Id)
	if err != nil {
		return err
	}

	p.CommentCount, err = Repos.Post.GetCommentCount(p.Id)
	if err != nil {
		return err
	}

	p.IsLiked, err = Repos.Post.GetUserReaction(userId, p.Id)
	if err != nil {
		return err
	}

	p.Categories, err = Repos.Category.GetNamesByPost(p.Id)
	return err
}

func enrichComment(c repository.Comment, userId int) (dblayer.Comment, error) {
	var result dblayer.Comment

	user, err := Repos.User.GetByID(c.UserID)
	if err != nil {
		return result, err
	}

	reaction, err := Repos.Reaction.GetReactionComment(c.ID, userId)
	if err != nil {
		return result, err
	}

	result = dblayer.Comment{
		Id:           c.ID,
		UserId:       c.UserID,
		Nickname:     user.Nickname,
		Created_at:   c.CreatedAt,
		Text:         c.Text,
		TimeAgo:      utilities.TimeAgo(c.CreatedAt),
		LikeCount:    reaction.Likes,
		DislikeCount: reaction.Dislikes,
		IsLiked:      reaction.IsLiked,
	}

	return result, nil
}

func enrichPostWithComments(p *dblayer.Post, userId int) error {
	if err := enrichPost(p, userId); err != nil {
		return err
	}

	comments, err := Repos.Comment.GetCommentsByPostPaginated(p.Id, 30, 0)
	if err != nil {
		return err
	}

	var enriched []dblayer.Comment
	for _, comment := range comments {
		c, err := enrichComment(comment, userId)
		if err != nil {
			return err
		}
		enriched = append(enriched, c)
	}
	p.Comments = enriched
	return nil
}

// =========================
// CREATE POST
// =========================

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/posts/create" {
		utilities.WriteJSON(w, http.StatusNotFound, "page not found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid form data", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	privacy := strings.TrimSpace(r.FormValue("privacy"))
	categories := r.MultipartForm.Value["categories"]
	allowedUsersRaw := strings.TrimSpace(r.FormValue("allowed_user_ids"))

	if title == "" || text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "title and text cannot be empty", nil)
		return
	}

	if len(title) > 255 || len(text) > 2000 {
		utilities.WriteJSON(w, http.StatusBadRequest, "title or text too long", nil)
		return
	}

	if len(categories) == 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "at least one category required", nil)
		return
	}

	if privacy == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "privacy is required", nil)
		return
	}

	if privacy != "public" && privacy != "almost_private" && privacy != "private" {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid privacy value", nil)
		return
	}

	categoryIDs, err := Repos.Category.GetIDsByNames(categories)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var imagePath string
	if file, header, err := r.FormFile("image"); err == nil {
		defer file.Close()
		if saved, saveErr := utilities.SaveImage(file, header, "uploads/posts/"); saveErr == nil {
			imagePath = saved
		} else {
			log.Printf("[CREATE POST] failed saving image: %v", saveErr)
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not save image", nil)
			return
		}
	}

	var allowedUserIDs []int
	if privacy == "private" && allowedUsersRaw != "" {
		parts := strings.Split(allowedUsersRaw, ",")
		for _, item := range parts {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			uid, parseErr := strconv.Atoi(item)
			if parseErr != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid allowed_user_ids list", nil)
				return
			}
			allowedUserIDs = append(allowedUserIDs, uid)
		}
	}

	post := &repository.Post{
		UserID:    userID,
		CreatedAt: time.Now(),
		Title:     title,
		Text:      text,
		Image:     imagePath,
		Privacy:   privacy,
	}

	if err := Repos.Post.CreatePost(post, allowedUserIDs, categoryIDs); err != nil {
		log.Printf("[CREATE POST] failed to create post: %v", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not create post", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusCreated, "post created successfully", map[string]any{
		"post_id": post.ID,
	})
}

// =========================
// POST RESOLVER
// =========================

func PostResolver(w http.ResponseWriter, r *http.Request) {
	segments := getPathSegments(r)
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "posts" {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	postID, err := strconv.Atoi(segments[2])
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
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

	if endpoint == "" {
		if r.Method == http.MethodGet {
			GetPostById(w, r)
			return
		}
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
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

		if err := Repos.Reaction.SetPostReaction(userID, postID, isLike); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		reactionCounts, err := Repos.Reaction.GetReactionPost(postID, userID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch reactions", nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "reaction updated", map[string]any{
			"post_id":  postID,
			"likes":    reactionCounts.Likes,
			"dislikes": reactionCounts.Dislikes,
			"is_liked": reactionCounts.IsLiked,
		})
	case "delete":
		if r.Method != http.MethodDelete {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if err := Repos.Post.DeletePost(postID, userID); err != nil {
			utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "deleted", map[string]any{"post_id": postID})
	default:
		utilities.WriteJSON(w, http.StatusNotFound, "unknown endpoint", nil)
	}
}

// =========================
// GET POSTS
// =========================
func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid query parameters", nil)
		return
	}

	categories := r.Form["categories"]
	liked := r.FormValue("liked_by_me") == "true"
	byMe := r.FormValue("posted_by_me") == "true"

	limit := 30
	lastID := 0

	if l := r.FormValue("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if id := r.FormValue("last_id"); id != "" {
		if v, err := strconv.Atoi(id); err == nil && v > 0 {
			lastID = v
		}
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	posts, err := Repos.Post.GetFilteredPosts(userID, categories, liked, byMe, limit, lastID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch posts", nil)
		return
	}

	enrichedPosts := make([]dblayer.Post, 0, len(posts))
	for _, post := range posts {
		dbPost := dblayer.Post{
			Id:         post.ID,
			UserId:     post.UserID,
			Created_at: post.CreatedAt,
			Title:      post.Title,
			Text:       post.Text,
			Image:      post.Image,
		}
		if err := enrichPost(&dbPost, userID); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to enrich posts", nil)
			return
		}
		enrichedPosts = append(enrichedPosts, dbPost)
	}

	utilities.WriteJSON(w, http.StatusOK, "posts fetched", enrichedPosts)
}

// =========================
// SINGLE POST
// =========================

func GetPostById(w http.ResponseWriter, r *http.Request) {
	segments := getPathSegments(r)
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "posts" {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	postID, err := strconv.Atoi(segments[2])
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	post, err := Repos.Post.GetVisiblePostByID(userID, postID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "post not found or not visible", nil)
		return
	}

	dbPost := dblayer.Post{
		Id:         post.ID,
		UserId:     post.UserID,
		Created_at: post.CreatedAt,
		Title:      post.Title,
		Text:       post.Text,
		Image:      post.Image,
	}
	if err := enrichPostWithComments(&dbPost, userID); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to enrich post", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "post fetched", dbPost)
}
