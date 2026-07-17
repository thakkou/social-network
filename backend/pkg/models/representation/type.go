package repModal

import "time"

type PostResponse struct {
	// We embed the main Post fields directly
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	TimeAgo   string    `json:"time_ago"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	Privacy   string    `json:"privacy"`

	// Aggregated counts for the post
	LikeCount    int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
	IsLiked      int `json:"is_liked"` // 1: liked, 0: none, -1: disliked

	// Associated Categories
	Categories []string `json:"categories"`

	// Nested Comments list
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	post_id   int       `db:"post_id"`
	CreatedAt time.Time `db:"created_at"`
	Title     string    `db:"title"`

	LikeCount    int `db:"like_count"`
	DislikeCount int `db:"dislike_count"`
	IsLiked      int `db:"is_liked"`
}

type POST_detaille struct{}

type ProfileResponse struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AboutMe   string `json:"aboutme"`

	IsPrivate int `json:"is_private"`

	FollowingStatus string `json:"following_status"`

	Followers []UserFollow   `json:"followers"`
	Following []UserFollow   `json:"following"`
	Posts     []PostResponse `json:"posts"`
}

type UserFollow struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}
