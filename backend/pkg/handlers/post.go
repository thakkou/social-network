package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "01social/pkg/db/sqlite"
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

var (
	postRepo     *repository.PostRepository
	categoryRepo *repository.CategoryRepository
)

func InitRepositories() {
	postRepo = repository.NewPostRepository(db.Database)
	categoryRepo = repository.NewCategoryRepository(db.Database)

	log.Println("Repositories initialized")
}

// =========================
// CORE POST ENRICHMENT
// =========================

func enrichPost(p *dblayer.Post, userId int) error {
	p.TimeAgo = utilities.TimeAgo(p.Created_at)

	// USER INFO
	// NOTE: nickname lookup isn't part of PostRepository/CategoryRepository
	// (it's user data, not post data). Left as a direct query for now; move
	// this into a UserRepository if/when one exists.
	if err := db.Database.QueryRow(
		"SELECT nickname FROM users WHERE id = ?",
		p.UserId,
	).Scan(&p.Nickname); err != nil {
		return err
	}

	// REACTIONS COUNT
	var err error
	p.LikeCount, p.DislikeCount, err = postRepo.GetReactionCounts(p.Id)
	if err != nil {
		return err
	}

	// USER REACTION
	p.IsLiked, err = postRepo.GetUserReaction(userId, p.Id)
	if err != nil {
		return err
	}

	// CATEGORIES
	p.Categories, err = categoryRepo.GetNamesByPost(p.Id)
	return err
}

func enrichPostWithComments(p *dblayer.Post, userId int) error {
	// if err := enrichPost(p, userId); err != nil {
	// 	return err
	// }

	// comments, err := postRepo.GetCommentsByPost(p.Id)
	// if err != nil {
	// 	return err
	// }

	// p.Comments = toDBLayerComments(comments)
	return nil
}

// =========================
// CREATE POST
// =========================

func CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start creating a posts")

	if r.URL.Path != "/api/posts/create" {
		utilities.WriteJSON(w, http.StatusNotFound, "Page Not Found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "Invalid form data", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, 403, "ononon", nil)
		return
	}

	title := r.FormValue("title")
	text := r.FormValue("text")
	privacy := r.FormValue("privacy")
	categories := r.MultipartForm.Value["categories"]

	// Debug incoming post data
	log.Println("========== CREATE POST ==========")
	log.Printf("UserID: %d", userID)
	log.Printf("Title: %q", title)
	log.Printf("Text Length: %d", len(text))
	log.Printf("Privacy: %q", privacy)
	log.Printf("Categories: %v", categories)

	if title == "" || text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "Title and text cannot be empty", nil)
		return
	}

	if len(title) > 255 || len(text) > 1000 {
		utilities.WriteJSON(w, http.StatusBadRequest, "Title or text too long", nil)
		return
	}

	if len(categories) == 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "At least one category required", nil)
		return
	}

	if privacy == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "Privacy is required", nil)
		return
	}

	categoryIDs, err := categoryRepo.GetIDsByNames(categories)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	log.Printf("Category IDs: %v", categoryIDs)

	// Optional image upload
	var imagePath string
	hasImage := false

	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		hasImage = true
		log.Println("image uploaded")

		path := "/uploads/image.png"
		imagePath = path
	} else {
		log.Println("no image uploaded")
	}

	// Final post data check
	log.Println("----- Final Post Data -----")
	log.Printf("UserID      : %d", userID)
	log.Printf("Title       : %s", title)
	log.Printf("Privacy     : %s", privacy)
	log.Printf("Has Image   : %t", hasImage)
	log.Printf("Image Path  : %s", imagePath)
	log.Printf("Categories  : %v", categories)
	log.Printf("Category IDs: %v", categoryIDs)
	log.Println("============================")

	fmt.Println("image path", imagePath, userID, categoryIDs)

	utilities.WriteJSON(w, http.StatusOK, "creat handleres3", nil)
}

// =========================
// POST RESOLVER
// =========================

func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	userId, _ := utilities.GetUserIDFromCookie(cookie.Value)
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	switch endpoint {

	// =========================
	// LIKE
	// =========================
	case "like":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if status, err := ReactToPost(userId, postId, 1); err != nil {
			utilities.WriteJSON(w, status, err.Error(), nil)
			return
		}

		likes, dislikes, _ := postRepo.GetReactionCounts(postId)
		reaction, _ := postRepo.GetUserReaction(userId, postId)

		utilities.WriteJSON(w, 200, "liked", map[string]any{
			"postId":   postId,
			"likes":    likes,
			"dislikes": dislikes,
			"isLike":   reaction, // 1, 0, or -1
		})

	// =========================
	// DISLIKE
	// =========================
	case "dislike":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if status, err := ReactToPost(userId, postId, -1); err != nil {
			utilities.WriteJSON(w, status, err.Error(), nil)
			return
		}

		likes, dislikes, _ := postRepo.GetReactionCounts(postId)
		reaction, _ := postRepo.GetUserReaction(userId, postId)

		utilities.WriteJSON(w, 200, "disliked", map[string]any{
			"postId":   postId,
			"likes":    likes,
			"dislikes": dislikes,
			"isLike":   reaction, // 1, 0, or -1
		})

	// =========================
	// DELETE
	// =========================
	case "delete":
		if r.Method != http.MethodDelete {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if err := postRepo.DeletePost(postId, userId); err != nil {
			utilities.WriteJSON(w, 403, err.Error(), nil)
			return
		}

		utilities.WriteJSON(w, 200, "deleted", map[string]any{
			"postId": postId,
		})

	default:
		utilities.WriteJSON(w, 404, "unknown endpoint", nil)
	}
}

// =========================
// GET POSTS
// =========================
func GetPosts(w http.ResponseWriter, r *http.Request) {
}

// later i will update it
// func GetPosts(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		utilities.WriteJSON(w, 405, "Method not allowed", nil)
// 		return
// 	}

// 	_ = r.ParseForm()

// 	categories := r.Form["categories"]
// 	liked := r.FormValue("my-liked-posts") == "true"
// 	byMe := r.FormValue("my-creat-posts") == "true"

// 	limit := 30
// 	lastID := 0

// 	if l := r.FormValue("limit"); l != "" {
// 		if v, err := strconv.Atoi(l); err == nil && v > 0 {
// 			limit = v
// 		}
// 	}
// 	if id := r.FormValue("lastId"); id != "" {
// 		if v, err := strconv.Atoi(id); err == nil && v > 0 {
// 			lastID = v
// 		}
// 	}

// 	var userID int
// 	if cookie, err := r.Cookie("session_id"); err == nil {
// 		userID, _ = utilities.GetUserIDFromCookie(cookie.Value)
// 	}

// 	posts, err := GetFilteredPosts(userID, categories, liked, byMe, limit, lastID)
// 	if err != nil {
// 		utilities.WriteJSON(w, 500, "error", nil)
// 		return
// 	}

// 	utilities.WriteJSON(w, 200, "ok", posts)
// }

// =========================
// SINGLE POST
// =========================

func GetPostById(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, 400, "Invalid ID", nil)
		return
	}

	limit := 15
	lastID := 0

	if l := r.URL.Query().Get("commentLimit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if lid := r.URL.Query().Get("commentLastId"); lid != "" {
		if v, err := strconv.Atoi(lid); err == nil && v > 0 {
			lastID = v
		}
	}
	fmt.Println("li", limit, lastID)

	// post, err := GetPostBasic(id)
	// if err != nil {
	// 	utilities.WriteJSON(w, 404, "Not found", nil)
	// 	return
	// }

	// comments, _ := postRepo.GetCommentsByPostPaginated(id, limit, lastID)
	// post.Comments = toDBLayerComments(comments)

	utilities.WriteJSON(w, 200, "ok", nil)
}
