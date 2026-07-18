package dbModal

import "time"

type Post struct {
	Id                      int
	UserId                  int
	Nickname                string
	Created_at              time.Time
	TimeAgo                 string
	Title                   string
	Text                    string
	LikeCount, DislikeCount int
	CommentCount            int
	IsLiked                 int // 1:liked, 0:none, -1:disliked
	Comments                []Comment
	Categories              []string
	Image                   string
}

type Comment struct {
	Id                      int
	UserId                  int
	Nickname                string
	Created_at              time.Time
	TimeAgo                 string
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:none, -1:disliked
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
