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
	"01social/pkg/ws"
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
	p.Firstname = user.Firstname
	p.Lastname = user.Lastname
	p.Avatar = user.Avatar

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
		Image:        c.Image,
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

// CreatePost creates a new post with optional image upload.
// @Summary Create a new post
// @Description Creates a new post with title, text, categories, privacy setting, optional image, and allowed user IDs for private posts.
// @Tags Posts
// @Accept mpfd
// @Produce json
// @Param title formData string true "Post title"
// @Param text formData string true "Post content"
// @Param privacy formData string true "Privacy level: public, almost_private, or private" Enums(public, almost_private, private)
// @Param categories formData []string false "Category names"
// @Param image formData file false "Optional post image"
// @Param allowed_user_ids formData string false "Comma-separated user IDs for private posts"
// @Success 201 {object} map[string]any "Post created successfully"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/posts/create [post]
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
		CreatedAt: time.Now().UTC(),
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

	if broadcaster, err := Repos.User.GetByID(userID); err == nil {
		ws.BroadcastExcept(strconv.Itoa(userID), "new_posts", map[string]any{
			"post_id":  post.ID,
			"user_id":  userID,
			"title":    title,
			"nickname": broadcaster.Nickname,
		})
	}

	utilities.WriteJSON(w, http.StatusCreated, "post created successfully", map[string]any{
		"post_id": post.ID,
	})
}

// =========================
// POST RESOLVER
// =========================

// PostResolver routes to the appropriate post action based on the URL path.
// GET  /api/posts/{id}          - Get single post with comments
// POST /api/posts/{id}/like     - Like a post
// POST /api/posts/{id}/dislike  - Dislike a post
// DELETE /api/posts/{id}/delete - Delete own post
// @Summary Post operations router
// @Description Handles single post fetch, like/dislike, and delete operations.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} dblayer.Post "Post data or reaction result"
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/posts/{id} [get]
// @Router /api/posts/{id}/like [post]
// @Router /api/posts/{id}/dislike [post]
// @Router /api/posts/{id}/delete [delete]
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

		// Notify the post author about the reaction
		postAuthor, err := Repos.Post.GetPostAuthor(postID)
		if err == nil && postAuthor != userID {
			ws.NotifyUser(strconv.Itoa(postAuthor), "like_posts", map[string]any{
				"post_id":  postID,
				"user_id":  userID,
				"reaction": endpoint,
			})

			Repos.Notification.Create(&repository.Notification{
				UserID:     postAuthor,
				ActorID:    userID,
				Type:       "post_reaction",
				ObjectType: "post",
				ObjectID:   postID,
			})
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
// GetPosts retrieves a filtered feed of posts.
// @Summary Get post feed
// @Description Returns a paginated, filterable feed of posts. Can filter by categories, liked by me, posted by me.
// @Tags Posts
// @Produce json
// @Param categories query []string false "Filter by category names" collectionFormat(multi)
// @Param liked_by_me query boolean false "Only posts liked by me"
// @Param posted_by_me query boolean false "Only my posts"
// @Param limit query integer false "Number of posts (default 30)" minimum(1)
// @Param last_id query integer false "Last post ID for pagination" minimum(0)
// @Success 200 {array} dblayer.Post "Posts fetched successfully"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid query parameters", nil)
		return
	}

	// Parse categories: supports both multiple values and comma-separated
	rawCategories := r.Form["categories"]
	var categories []string
	for _, c := range rawCategories {
		for _, part := range strings.Split(c, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				categories = append(categories, part)
			}
		}
	}
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
		Privacy:    post.Privacy,
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
		Privacy:    post.Privacy,
	}
	if err := enrichPostWithComments(&dbPost, userID); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to enrich post", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "post fetched", dbPost)
}
