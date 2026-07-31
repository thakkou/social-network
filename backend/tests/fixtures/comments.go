// Package fixtures provides typed constants and helpers for the seeder_test.sql data.
//
// NOTE: CommentsByPost only covers the first ~20 of 88 seeded comments (posts 1-22).
// If your test needs comments on later posts, add them here or use direct DB queries.
package fixtures

// CommentInfo holds reference info for a seeded comment.
type CommentInfo struct {
	ID     int
	UserID int
	PostID int
	Snippet string
}

// CommentsByPost maps post IDs to their seeded comments.
var CommentsByPost = map[int][]CommentInfo{
	1: {
		{ID: 1, UserID: Bob.ID, PostID: 1, Snippet: "Welcome! Great to have you here."},
		{ID: 2, UserID: David.ID, PostID: 1, Snippet: "Nice first post"},
	},
	2:   {{ID: 3, UserID: Alice.ID, PostID: 2, Snippet: "Looks amazing"}},
	4:   {{ID: 4, UserID: Farid.ID, PostID: 4, Snippet: "Solid tip"}},
	5:   {{ID: 5, UserID: Alice.ID, PostID: 5, Snippet: "Thanks for sharing"}},
	6:   {{ID: 7, UserID: Alice.ID, PostID: 6, Snippet: "productive morning"}, {ID: 8, UserID: Bob.ID, PostID: 6, Snippet: "Coffee first"}},
	7:   {{ID: 9, UserID: David.ID, PostID: 7, Snippet: "love to see"}, {ID: 10, UserID: Farid.ID, PostID: 7, Snippet: "Sunsets"}},
	8:   {{ID: 11, UserID: Alice.ID, PostID: 8, Snippet: "Hope weather"}, {ID: 12, UserID: Hugo.ID, PostID: 8, Snippet: "Enjoy hike"}},
	9:   {{ID: 13, UserID: Bob.ID, PostID: 9, Snippet: "Go interfaces"}, {ID: 14, UserID: Farid.ID, PostID: 9, Snippet: "generics"}},
	11:  {{ID: 6, UserID: Bob.ID, PostID: 11, Snippet: "hear it"}},
	12:  {{ID: 15, UserID: Bob.ID, PostID: 12, Snippet: "masterpiece"}, {ID: 16, UserID: Farid.ID, PostID: 12, Snippet: "re-read"}},
	13:  {{ID: 17, UserID: David.ID, PostID: 13, Snippet: "tough"}, {ID: 18, UserID: Alice.ID, PostID: 13, Snippet: "Stick with it"}},
	15:  {{ID: 19, UserID: Farid.ID, PostID: 15, Snippet: "Consistency"}},
	16:  {{ID: 20, UserID: Alice.ID, PostID: 16, Snippet: "DAW"}},
	17:  {{ID: 21, UserID: Bob.ID, PostID: 17, Snippet: "Count me in"}},
	18:  {{ID: 22, UserID: Grace.ID, PostID: 18, Snippet: "stunning"}},
}

// CommentByID returns a comment from any post, or nil.
func CommentByID(commentID int) *CommentInfo {
	for _, comments := range CommentsByPost {
		for _, c := range comments {
			if c.ID == commentID {
				return &c
			}
		}
	}
	return nil
}

// CommentsForPost returns all seeded comments for a post.
func CommentsForPost(postID int) []CommentInfo {
	return CommentsByPost[postID]
}
