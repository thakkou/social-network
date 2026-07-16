package models

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
