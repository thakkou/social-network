package seeder

import (
	"database/sql"
	"log"
	"time"
)

func seedPosts(db *sql.DB, u []int) ([]int, error) {
	type p struct {
		UserID      int
		Title, Text string
		Image       string
		Privacy     string
	}
	posts := []p{
		// Original posts
		{u[0], "Hello world", "My very first post on this network!", "/uploads/seeder/posts/dev.jpeg", "public"},
		{u[1], "Weekend trip", "Just got back from the mountains 🏔️", "/uploads/seeder/posts/travel.jpeg", "public"},
		{u[3], "", "Debugging is 90% of my job today.", "/uploads/seeder/posts/dev2.jpeg", "almost_private"},
		{u[5], "Go tip", "context.Context should be your first param, always.", "/uploads/seeder/posts/dev3.jpeg", "public"},
		{u[1], "Private thoughts", "Only a few people should see this.", "/uploads/seeder/posts/learning.jpeg", "private"},
		{u[1], "Morning routine", "Coffee, reading, then coding.", "/uploads/seeder/posts/cooking1.jpeg", "public"},
		{u[1], "Photography", "Took some beautiful sunset photos today.", "/uploads/seeder/posts/travel2.jpeg", "public"},
		{u[1], "Weekend plans", "Thinking about hiking this weekend.", "/uploads/seeder/posts/cooking2.jpeg", "almost_private"},
		{u[2], "Learning Go", "Interfaces finally clicked today!", "/uploads/seeder/posts/learning2.jpeg", "public"},
		{u[4], "New camera", "Testing out my new lens today.", "/uploads/seeder/posts/travel.jpeg", "almost_private"},
		{u[7], "New track", "Dropping a new beat this Friday 🎵", "", "public"},

		// More posts for richer feed
		{u[0], "Distributed Systems", "Reading 'Designing Data-Intensive Applications'. Highly recommend!", "/uploads/seeder/posts/dev3.jpeg", "public"},
		{u[5], "Rust vs Go", "Trying out Rust for a CLI tool. The borrow checker is something else.", "/uploads/seeder/posts/dev2.jpeg", "public"},
		{u[3], "Weekend project", "Built a tiny load balancer in Go over the weekend.", "", "almost_private"},
		{u[6], "New hobby", "Started learning the piano! Any tips for beginners?", "", "public"},
		{u[7], "Studio update", "New soundproof panels in the studio. Check it out!", "/uploads/seeder/posts/dev.jpeg", "public"},
		{u[2], "Book club", "Reading 'The Pragmatic Programmer' with the Go group. Join us!", "", "public"},
		{u[4], "Sunset shots", "Golden hour at the lake today. Nature is healing.", "/uploads/seeder/posts/travel2.jpeg", "public"},

		// New seed posts for even richer data
		{u[8], "Art Exhibition", "Visited the modern art museum today. Absolutely inspiring pieces!", "/uploads/seeder/posts/travel.jpeg", "public"},
		{u[9], "Open Source Saturday", "Contributed to my first open source project today! Feeling great.", "/uploads/seeder/posts/dev2.jpeg", "public"},
		{u[6], "Baking bread", "Made sourdough from scratch for the first time. Turned out amazing!", "/uploads/seeder/posts/cooking1.jpeg", "public"},
		{u[0], "Microservices talk", "Went to a great meetup about microservices patterns. Event sourcing is fascinating.", "/uploads/seeder/posts/dev3.jpeg", "public"},
		{u[2], "Coding playlist", "Share your favorite coding playlist! I need new recommendations.", "", "public"},
		{u[5], "Docker tips", "Multi-stage builds are a game changer for reducing image size. Here's how...", "/uploads/seeder/posts/dev.jpeg", "public"},
		{u[3], "Night sky", "Captured the Milky Way with my telescope last night. Astro-photography is hard!", "/uploads/seeder/posts/travel2.jpeg", "almost_private"},
		{u[7], "New single out now", "My new single 'Midnight Code' is out on all platforms! Link in bio 🎶", "", "public"},
		{u[8], "Digital art", "Experimenting with procedural art generation using Go. Unexpectedly beautiful results!", "/uploads/seeder/posts/learning.jpeg", "public"},
		{u[1], "Travel tips", "Best travel destinations for digital nomads in 2026: full guide in thread.", "/uploads/seeder/posts/travel.jpeg", "public"},
		{u[4], "Film photography", "Got my first film roll developed. There's something magical about analog.", "/uploads/seeder/posts/learning2.jpeg", "public"},
		{u[9], "Terminal tools", "My favorite terminal tools: fzf, ripgrep, bat, and lazygit. What are yours?", "/uploads/seeder/posts/dev3.jpeg", "public"},
		{u[6], "Hiking adventure", "Summited Mount Toubkal! Toughest hike of my life but the view was worth it.", "/uploads/seeder/posts/travel2.jpeg", "public"},
		{u[0], "AI pair programming", "Been using AI tools for code reviews. They catch things I miss!", "/uploads/seeder/posts/dev2.jpeg", "public"},
		{u[3], "Minimalist setup", "My new WFH setup: minimal, clean, productive. Less is more.", "/uploads/seeder/posts/dev.jpeg", "public"},
	}

	query := `INSERT INTO POSTS (user_id, created_at, title, text, image, privacy) VALUES (?, ?, ?, ?, ?, ?)`
	ids := make([]int, 0, len(posts))
	now := time.Now()
	for i, post := range posts {
		title := sql.NullString{String: post.Title, Valid: post.Title != ""}
		createdAt := now.Add(-time.Duration(len(posts)-i) * time.Hour).Format("2006-01-02 15:04:05")
		image := sql.NullString{String: post.Image, Valid: post.Image != ""}
		res, err := db.Exec(query, post.UserID, createdAt, title, post.Text, image, post.Privacy)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] posts: %d\n", len(ids))
	return ids, nil
}

func seedPostAllowedUsers(db *sql.DB, postIDs, u []int) error {
	privatePostID := postIDs[4]
	allowed := []int{u[6], u[0]}
	query := `INSERT INTO POST_ALLOWED_USERS (post_id, user_id) VALUES (?, ?)`
	for _, uid := range allowed {
		if _, err := db.Exec(query, privatePostID, uid); err != nil {
			return err
		}
	}
	log.Printf("[SEED] post_allowed_users: %d\n", len(allowed))
	return nil
}

func seedPostCategories(db *sql.DB, postIDs []int) error {
	catID := func(name string) (int, error) {
		var id int
		err := db.QueryRow("SELECT id FROM CATEGORY WHERE name = ?", name).Scan(&id)
		return id, err
	}
	general, err := catID("General")
	if err != nil {
		return err
	}
	travel, err := catID("Travel")
	if err != nil {
		return err
	}
	edu, err := catID("Education")
	if err != nil {
		return err
	}
	ent, err := catID("Entertainment")
	if err != nil {
		return err
	}
	sports, err := catID("Sports")
	if err != nil {
		return err
	}
	links := map[int][]int{
		postIDs[0]:  {general},
		postIDs[1]:  {travel},
		postIDs[3]:  {edu},
		postIDs[6]:  {ent},
		postIDs[11]: {edu},
		postIDs[12]: {edu},
		postIDs[18]: {ent},
		postIDs[20]: {general},
		postIDs[21]: {edu},
		postIDs[23]: {edu},
		postIDs[24]: {edu},
		postIDs[25]: {ent},
		postIDs[26]: {travel},
		postIDs[27]: {travel, sports},
		postIDs[28]: {ent},
	}
	query := `INSERT INTO POST_CATEGORY (post_id, category_id) VALUES (?, ?)`
	count := 0
	for postID, cats := range links {
		for _, c := range cats {
			if _, err := db.Exec(query, postID, c); err != nil {
				return err
			}
			count++
		}
	}
	log.Printf("[SEED] post_category: %d\n", count)
	return nil
}
