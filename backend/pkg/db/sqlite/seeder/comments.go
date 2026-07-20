package seeder

import (
	"log"
	"time"

	"database/sql"
)

func seedComments(db *sql.DB, u, postIDs []int) ([]int, error) {
	type c struct {
		UserID, PostID int
		Text           string
	}
	comments := []c{
		{u[1], postIDs[0], "Welcome! Great to have you here."},
		{u[3], postIDs[0], "Nice first post 🎉"},
		{u[0], postIDs[1], "Looks amazing, take me next time!"},
		{u[5], postIDs[3], "Solid tip, saved me a bug last week."},
		{u[0], postIDs[4], "Thanks for sharing this with me."},
		{u[1], postIDs[10], "Can't wait to hear it!"},
		{u[0], postIDs[5], "That sounds like a productive morning!"},
		{u[1], postIDs[5], "Coffee first is always the right choice ☕"},
		{u[3], postIDs[6], "Would love to see those photos!"},
		{u[5], postIDs[6], "Sunsets are the best 🌅"},
		{u[0], postIDs[7], "Hope the weather stays nice!"},
		{u[7], postIDs[7], "Enjoy your hike!"},
		{u[1], postIDs[8], "Go interfaces are confusing at first 😄"},
		{u[5], postIDs[8], "Wait until you discover generics!"},

		// More comments on new posts
		{u[1], postIDs[11], "That book is a masterpiece! Chapter 5 on replication is gold."},
		{u[5], postIDs[11], "I should re-read this. Great recommendation."},
		{u[3], postIDs[12], "Borrow checker is tough but rewarding."},
		{u[0], postIDs[12], "Stick with it! Rust's ecosystem is growing fast."},
		{u[5], postIDs[14], "Consistency is key! 10 minutes a day makes a big difference."},
		{u[0], postIDs[15], "Nice setup! What DAW are you using?"},
		{u[1], postIDs[16], "Count me in! I'm halfway through the book."},
		{u[6], postIDs[17], "Absolutely stunning 😍"},
	}

	query := `INSERT INTO COMMENTS (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)`
	ids := make([]int, 0, len(comments))
	now := time.Now()
	for i, cm := range comments {
		createdAt := now.Add(-time.Duration(len(comments)-i) * time.Minute * 30).Format("2006-01-02 15:04:05")
		res, err := db.Exec(query, cm.UserID, cm.PostID, createdAt, cm.Text)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] comments: %d\n", len(ids))
	return ids, nil
}

func seedReactions(db *sql.DB, u, postIDs, commentIDs []int) error {
	type pr struct {
		UserID, PostID, IsLike int
	}
	postReactions := []pr{
		{u[1], postIDs[0], 1},
		{u[1], postIDs[8], 1},
		{u[1], postIDs[7], 1},
		{u[3], postIDs[0], 1},
		{u[5], postIDs[1], 1},
		{u[0], postIDs[3], 1},
		{u[7], postIDs[3], -1},
		{u[0], postIDs[6], 1},

		// More reactions
		{u[5], postIDs[0], 1},
		{u[0], postIDs[10], 1},
		{u[3], postIDs[11], 1},
		{u[1], postIDs[12], -1},
		{u[7], postIDs[14], 1},
		{u[2], postIDs[16], 1},
	}
	q1 := `INSERT INTO POST_REACTIONS (user_id, post_id, is_like) VALUES (?, ?, ?)`
	for _, r := range postReactions {
		if _, err := db.Exec(q1, r.UserID, r.PostID, r.IsLike); err != nil {
			return err
		}
	}

	type cr struct {
		UserID, CommentID, IsLike int
	}
	commentReactions := []cr{
		{u[0], commentIDs[0], 1},
		{u[5], commentIDs[3], 1},
		{u[1], commentIDs[4], 1},
	}
	q2 := `INSERT INTO COMMENT_REACTIONS (user_id, comment_id, is_like) VALUES (?, ?, ?)`
	for _, r := range commentReactions {
		if _, err := db.Exec(q2, r.UserID, r.CommentID, r.IsLike); err != nil {
			return err
		}
	}

	log.Printf("[SEED] post_reactions: %d, comment_reactions: %d\n", len(postReactions), len(commentReactions))
	return nil
}
