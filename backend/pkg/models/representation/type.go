package representation

type Post struct {
	ID    int
	Title string
}
type ProfileResponse struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AboutMe   string `json:"aboutme"`

	IsPrivate bool `json:"is_private"`

	FollowingStatus string `json:"following_status"`

	Followers int `json:"followers"`
	Following int `json:"following"`
	Posts     int `json:"posts"`

	PostsList []Post `json:"posts"`
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
