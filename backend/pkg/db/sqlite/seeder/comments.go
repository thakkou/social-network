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
		// Original comments
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

		// More comments on original posts
		{u[1], postIDs[11], "That book is a masterpiece! Chapter 5 on replication is gold."},
		{u[5], postIDs[11], "I should re-read this. Great recommendation."},
		{u[3], postIDs[12], "Borrow checker is tough but rewarding."},
		{u[0], postIDs[12], "Stick with it! Rust's ecosystem is growing fast."},
		{u[5], postIDs[14], "Consistency is key! 10 minutes a day makes a big difference."},
		{u[0], postIDs[15], "Nice setup! What DAW are you using?"},
		{u[1], postIDs[16], "Count me in! I'm halfway through the book."},
		{u[6], postIDs[17], "Absolutely stunning 😍"},

		// Comments on new seed posts
		{u[0], postIDs[18], "Modern art is so thought-provoking. Which exhibit was your favorite?"},
		{u[5], postIDs[18], "Wish I could visit! Send pictures."},
		{u[1], postIDs[19], "That's awesome! Which project did you contribute to?"},
		{u[7], postIDs[19], "Open source is the way to go! Keep it up 💪"},
		{u[3], postIDs[20], "That looks delicious! Drop the recipe 🍞"},
		{u[8], postIDs[20], "Sourdough is the best. Try adding rosemary next time!"},
		{u[5], postIDs[21], "Event sourcing is great for audit logs. Check out EventStore!"},
		{u[1], postIDs[21], "I was at that meetup too! Great talk."},
		{u[9], postIDs[22], "Lo-fi hip hop beats to code to, always."},
		{u[0], postIDs[22], "I like ambient electronic. Brian Eno is great."},
		{u[3], postIDs[23], "Multi-stage builds saved me 60% image size. Great tip!"},
		{u[7], postIDs[23], "Don't forget to use .dockerignore too!"},
		{u[0], postIDs[24], "Absolutely breathtaking! What equipment did you use?"},
		{u[6], postIDs[24], "I tried astrophotography once, it's so difficult. Great shot!"},
		{u[1], postIDs[25], "Congrats on the release! Adding to my playlist now 🎧"},
		{u[5], postIDs[25], "Love the title! 'Midnight Code' hits different."},
		{u[3], postIDs[26], "Procedural generation is so cool. What library did you use?"},
		{u[9], postIDs[27], "Portugal is amazing for digital nomads. Highly recommend!"},
		{u[0], postIDs[27], "Thailand is also incredible. Great community there."},
		{u[1], postIDs[28], "Analog has such a unique feel. What camera are you using?"},
		{u[7], postIDs[29], "ripgrep is a lifesaver! I also love jq for JSON processing."},
		{u[3], postIDs[29], "lazygit changed my workflow completely."},
		{u[5], postIDs[30], "Wow that's impressive! How long did the hike take?"},
		{u[8], postIDs[30], "Toubkal is on my bucket list! Any tips for beginners?"},
		{u[1], postIDs[31], "AI code review is surprisingly good. Which tool do you use?"},
		{u[6], postIDs[31], "Just be careful about sensitive data in the prompts!"},
		{u[0], postIDs[32], "Love the clean setup! What monitor is that?"},
		{u[5], postIDs[32], "Minimalist setups are the best for focus."},
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
		// Original reactions
		{u[1], postIDs[0], 1},
		{u[1], postIDs[8], 1},
		{u[1], postIDs[7], 1},
		{u[3], postIDs[0], 1},
		{u[5], postIDs[1], 1},
		{u[0], postIDs[3], 1},
		{u[7], postIDs[3], -1},
		{u[0], postIDs[6], 1},
		{u[5], postIDs[0], 1},
		{u[0], postIDs[10], 1},
		{u[3], postIDs[11], 1},
		{u[1], postIDs[12], -1},
		{u[7], postIDs[14], 1},
		{u[2], postIDs[16], 1},

		// More reactions on new posts
		{u[0], postIDs[18], 1},
		{u[5], postIDs[18], 1},
		{u[1], postIDs[19], 1},
		{u[7], postIDs[19], 1},
		{u[3], postIDs[20], 1},
		{u[8], postIDs[20], -1},
		{u[5], postIDs[21], 1},
		{u[1], postIDs[21], 1},
		{u[0], postIDs[22], 1},
		{u[9], postIDs[22], 1},
		{u[3], postIDs[23], 1},
		{u[7], postIDs[23], 1},
		{u[0], postIDs[24], 1},
		{u[6], postIDs[24], 1},
		{u[1], postIDs[25], 1},
		{u[5], postIDs[25], 1},
		{u[3], postIDs[26], 1},
		{u[0], postIDs[27], 1},
		{u[1], postIDs[28], 1},
		{u[7], postIDs[29], 1},
		{u[5], postIDs[30], 1},
		{u[8], postIDs[30], 1},
		{u[1], postIDs[31], 1},
		{u[6], postIDs[31], -1},
		{u[0], postIDs[32], 1},
		{u[5], postIDs[32], 1},
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
		// Original comment reactions
		{u[0], commentIDs[0], 1},
		{u[5], commentIDs[3], 1},
		{u[1], commentIDs[4], 1},

		// More comment reactions
		{u[0], commentIDs[5], 1},
		{u[3], commentIDs[6], 1},
		{u[5], commentIDs[8], -1},
		{u[1], commentIDs[10], 1},
		{u[0], commentIDs[12], 1},
		{u[7], commentIDs[14], 1},
		{u[1], commentIDs[15], 1},
		{u[5], commentIDs[18], 1},
		{u[0], commentIDs[20], 1},
		{u[3], commentIDs[22], 1},
		{u[5], commentIDs[24], 1},
		{u[7], commentIDs[26], 1},
		{u[1], commentIDs[28], 1},
		{u[0], commentIDs[30], 1},
		{u[5], commentIDs[32], 1},
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
