package fixtures

// PostInfo holds key details about a seeded post.
type PostInfo struct {
	ID       int
	AuthorID int
	Title    string
	Privacy  string // "public" | "almost_private" | "private"
	Category string // Primary category name
}

// PostsByAuthor maps author IDs to their seeded posts.
var PostsByAuthor = map[int][]PostInfo{
	Alice.ID: {
		{ID: 1, AuthorID: Alice.ID, Title: "Hello world", Privacy: "public", Category: "General"},
		{ID: 12, AuthorID: Alice.ID, Title: "Distributed Systems", Privacy: "public", Category: "Education"},
		{ID: 22, AuthorID: Alice.ID, Title: "Microservices talk", Privacy: "public", Category: "Education"},
		{ID: 32, AuthorID: Alice.ID, Title: "AI pair programming", Privacy: "public", Category: "Education"},
		{ID: 39, AuthorID: Alice.ID, Title: "Promotion news", Privacy: "public", Category: "Business"},
		{ID: 48, AuthorID: Alice.ID, Title: "Legacy code", Privacy: "public", Category: "Education"},
		{ID: 58, AuthorID: Alice.ID, Title: "SQLite vs PostgreSQL", Privacy: "public", Category: "Education"},
	},
	Bob.ID: {
		{ID: 2, AuthorID: Bob.ID, Title: "Weekend trip", Privacy: "public", Category: "Travel"},
		{ID: 5, AuthorID: Bob.ID, Title: "Private thoughts", Privacy: "private", Category: ""},
		{ID: 6, AuthorID: Bob.ID, Title: "Morning routine", Privacy: "public", Category: "Lifestyle"},
		{ID: 7, AuthorID: Bob.ID, Title: "Photography", Privacy: "public", Category: "Entertainment"},
		{ID: 8, AuthorID: Bob.ID, Title: "Weekend plans", Privacy: "almost_private", Category: "Sports"},
		{ID: 28, AuthorID: Bob.ID, Title: "Travel tips", Privacy: "public", Category: "Travel"},
		{ID: 42, AuthorID: Bob.ID, Title: "Sahara expedition", Privacy: "public", Category: "Travel"},
		{ID: 52, AuthorID: Bob.ID, Title: "Alps hiking", Privacy: "public", Category: "Sports"},
		{ID: 62, AuthorID: Bob.ID, Title: "Budget travel", Privacy: "public", Category: "Travel"},
	},
	Chloe.ID: {
		{ID: 9, AuthorID: Chloe.ID, Title: "Learning Go", Privacy: "public", Category: "Education"},
		{ID: 17, AuthorID: Chloe.ID, Title: "Book club", Privacy: "public", Category: "Education"},
		{ID: 23, AuthorID: Chloe.ID, Title: "Coding playlist", Privacy: "public", Category: "Entertainment"},
		{ID: 34, AuthorID: Chloe.ID, Title: "Pointers in C", Privacy: "public", Category: "Education"},
		{ID: 46, AuthorID: Chloe.ID, Title: "Clean Code review", Privacy: "public", Category: "Personal Dev"},
		{ID: 56, AuthorID: Chloe.ID, Title: "Linear algebra", Privacy: "public", Category: "Education"},
	},
	David.ID: {
		{ID: 3, AuthorID: David.ID, Title: "", Privacy: "almost_private", Category: "Education"},
		{ID: 14, AuthorID: David.ID, Title: "Weekend project", Privacy: "almost_private", Category: "Education"},
		{ID: 25, AuthorID: David.ID, Title: "Night sky", Privacy: "almost_private", Category: "Education"},
		{ID: 33, AuthorID: David.ID, Title: "Minimalist setup", Privacy: "public", Category: "Lifestyle"},
		{ID: 41, AuthorID: David.ID, Title: "Homemade sushi", Privacy: "public", Category: "Food & Cooking"},
		{ID: 51, AuthorID: David.ID, Title: "Weather app", Privacy: "public", Category: "Education"},
		{ID: 61, AuthorID: David.ID, Title: "Home server", Privacy: "public", Category: "Education"},
	},
	Emma.ID: {
		{ID: 10, AuthorID: Emma.ID, Title: "New camera", Privacy: "almost_private", Category: "Entertainment"},
		{ID: 18, AuthorID: Emma.ID, Title: "Sunset shots", Privacy: "public", Category: "Travel"},
		{ID: 29, AuthorID: Emma.ID, Title: "Film photography", Privacy: "public", Category: "Entertainment"},
		{ID: 38, AuthorID: Emma.ID, Title: "Street photography", Privacy: "public", Category: "Entertainment"},
		{ID: 50, AuthorID: Emma.ID, Title: "Tokyo guide", Privacy: "public", Category: "Travel"},
		{ID: 60, AuthorID: Emma.ID, Title: "Film photography", Privacy: "public", Category: "Culture"},
	},
	Farid.ID: {
		{ID: 4, AuthorID: Farid.ID, Title: "Go tip", Privacy: "public", Category: "Education"},
		{ID: 13, AuthorID: Farid.ID, Title: "Rust vs Go", Privacy: "public", Category: "Education"},
		{ID: 24, AuthorID: Farid.ID, Title: "Docker tips", Privacy: "public", Category: "Education"},
		{ID: 43, AuthorID: Farid.ID, Title: "Kubernetes thoughts", Privacy: "public", Category: "Education"},
		{ID: 53, AuthorID: Farid.ID, Title: "gRPC vs REST", Privacy: "public", Category: "Education"},
		{ID: 63, AuthorID: Farid.ID, Title: ".vimrc secrets", Privacy: "public", Category: "Education"},
	},
	Grace.ID: {
		{ID: 15, AuthorID: Grace.ID, Title: "New hobby", Privacy: "public", Category: "Personal Dev"},
		{ID: 21, AuthorID: Grace.ID, Title: "Baking bread", Privacy: "public", Category: "General"},
		{ID: 31, AuthorID: Grace.ID, Title: "Hiking adventure", Privacy: "public", Category: "Travel"},
		{ID: 35, AuthorID: Grace.ID, Title: "Morning yoga", Privacy: "public", Category: "Health & Fitness"},
		{ID: 45, AuthorID: Grace.ID, Title: "First marathon", Privacy: "public", Category: "Sports"},
		{ID: 54, AuthorID: Grace.ID, Title: "Vegan banana bread", Privacy: "public", Category: "Food & Cooking"},
	},
	Hugo.ID: {
		{ID: 11, AuthorID: Hugo.ID, Title: "New track", Privacy: "public", Category: "Entertainment"},
		{ID: 16, AuthorID: Hugo.ID, Title: "Studio update", Privacy: "public", Category: "Entertainment"},
		{ID: 26, AuthorID: Hugo.ID, Title: "New single out now", Privacy: "public", Category: "Entertainment"},
		{ID: 40, AuthorID: Hugo.ID, Title: "New EP done", Privacy: "public", Category: "Entertainment"},
		{ID: 49, AuthorID: Hugo.ID, Title: "Studio tour", Privacy: "public", Category: "Entertainment"},
		{ID: 59, AuthorID: Hugo.ID, Title: "Concert night", Privacy: "public", Category: "Entertainment"},
	},
	Bella.ID: {
		{ID: 19, AuthorID: Bella.ID, Title: "Art Exhibition", Privacy: "public", Category: "Entertainment"},
		{ID: 27, AuthorID: Bella.ID, Title: "Digital art", Privacy: "public", Category: "Travel"},
		{ID: 36, AuthorID: Bella.ID, Title: "New painting", Privacy: "public", Category: "Culture"},
		{ID: 47, AuthorID: Bella.ID, Title: "UI design tips", Privacy: "public", Category: "Education"},
		{ID: 57, AuthorID: Bella.ID, Title: "Forest photography", Privacy: "public", Category: "Travel"},
	},
	Jack.ID: {
		{ID: 20, AuthorID: Jack.ID, Title: "Open Source Saturday", Privacy: "public", Category: "General"},
		{ID: 30, AuthorID: Jack.ID, Title: "Terminal tools", Privacy: "public", Category: "Education"},
		{ID: 37, AuthorID: Jack.ID, Title: "Rust CLI tool", Privacy: "public", Category: "Education"},
		{ID: 44, AuthorID: Jack.ID, Title: "PR merged!", Privacy: "public", Category: "General"},
		{ID: 55, AuthorID: Jack.ID, Title: "Custom keyboard", Privacy: "public", Category: "Lifestyle"},
	},
}

// PostByID returns the PostInfo for the given ID by searching all authors.
func PostByID(postID int) *PostInfo {
	for _, posts := range PostsByAuthor {
		for _, p := range posts {
			if p.ID == postID {
				return &p
			}
		}
	}
	return nil
}

// PostsByAuthorID returns the seeded posts by a specific author ID.
func PostsByAuthorID(authorID int) []PostInfo {
	return PostsByAuthor[authorID]
}

// PostsByUser returns the seeded posts by a specific user (the preferred
// way over PostsByAuthorID since it uses the named UserInfo constant).
func PostsByUser(u UserInfo) []PostInfo {
	return PostsByAuthor[u.ID]
}

// PostIDsOf returns just the post IDs for a given author.
func PostIDsOf(authorID int) []int {
	posts := PostsByAuthor[authorID]
	ids := make([]int, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}
	return ids
}

// PublicPostIDs returns IDs of all public posts.
func PublicPostIDs() []int {
	var ids []int
	for _, posts := range PostsByAuthor {
		for _, p := range posts {
			if p.Privacy == "public" {
				ids = append(ids, p.ID)
			}
		}
	}
	return ids
}

// PrivatePostIDs returns IDs of private posts (post 5 - Bob's private post).
func PrivatePostIDs() []int {
	return []int{5}
}

// AlmostPrivatePostIDs returns IDs of almost_private posts.
func AlmostPrivatePostIDs() []int {
	return []int{3, 8, 10, 14, 25}
}
