package repModal

type Post struct {
	ID       int
	Title    string
	Text     string
	Image    string
	comment  int
	like     int
	dislike  int
	is_liked int
}

type ProfileResponse struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AboutMe   string `json:"aboutme"`

	IsPrivate int `json:"is_private"`

	FollowingStatus string `json:"following_status"`

	Followers []UserFollow `json:"followers"`
	Following []UserFollow `json:"following"`
	Posts     []Post       `json:"posts"`
}

type UserFollow struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}
