package dbModal

import "time"

type Post struct {
	Id           int       `json:"id"`
	UserId       int       `json:"user_id"`
	Nickname     string    `json:"nickname"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Avatar       string    `json:"avatar"`
	Created_at   time.Time `json:"created_at"`
	TimeAgo      string    `json:"time_ago"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	Image        string    `json:"image"`
	Privacy      string    `json:"privacy"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CommentCount int       `json:"comment_count"`
	IsLiked      int       `json:"is_liked"` // 1:liked, 0:none, -1:disliked
	Comments     []Comment `json:"comments"`
	Categories   []string  `json:"categories"`
}

type Comment struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	Nickname     string `json:"nickname"`
	Created_at   time.Time `json:"created_at"`
	TimeAgo      string `json:"time_ago"`
	Text         string `json:"text"`
	Image        string `json:"image"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	IsLiked      int    `json:"is_liked"` // 1:liked, 0:none, -1:disliked
}

type User struct {
	Id int

	// sanitize all data in frontend and auth.go
	Firstname string `json:"firstname"` // seed.go
	Lastname  string `json:"lastname"`  // seed.go
	Email     string `json:"email"`
	Password  string `json:"password"`
	BirthDate string `json:"birthDate"` // seed.go
	Nickname  string `json:"nickname"`  // opt
	AboutMe   string `json:"aboutme"`   // opt

	// Avatar type? `json:"avatar"` // opt

	// OAUTH
	// github + (google requires other apis => 'name')
	// Message string // (NOT STORED) ???
	// Login string `json:"login"` // oauth
	// Picture string `json:"picture"`    // gmail picture: sometimes cannot be loaded!
	// Avatar  string `json:"avatar_url"` // github avatar
}
