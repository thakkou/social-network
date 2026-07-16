// Package seed populates the database with sample data for local development
// and testing. It matches the schema in your migrations: USERS, FOLLOWS,
// POSTS, POST_ALLOWED_USERS, POST_CATEGORY, COMMENTS, POST_REACTIONS,
// COMMENT_REACTIONS, GROUPS, GROUP_MEMBERS, GROUP_INVITES, GROUP_REQUESTS,
// GROUP_MESSAGES, GROUP_EVENTS, EVENT_RESPONSES, GROUP_POSTS,
// GROUP_POST_COMMENTS, CONVERSATIONS, MESSAGES, NOTIFICATIONS.
//
// Usage (add a flag to your main.go, or call seed.Run directly):
//
//	go run . --seed
//
// or from code:
//
//	if err := seed.Run(db.Database); err != nil {
//	    log.Fatalf("seed failed: %v", err)
//	}
package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Run truncates the seedable tables and inserts a consistent set of sample
// data. Every seeded user's password is "Password123!".
func Run(db *sql.DB) error {
	log.Println("[SEED] starting...")

	if err := reset(db); err != nil {
		return fmt.Errorf("reset: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	pw := string(hash)

	userIDs, err := seedUsers(db, pw)
	if err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}

	if err := seedFollows(db, userIDs); err != nil {
		return fmt.Errorf("seedFollows: %w", err)
	}

	postIDs, err := seedPosts(db, userIDs)
	if err != nil {
		return fmt.Errorf("seedPosts: %w", err)
	}

	if err := seedPostAllowedUsers(db, postIDs, userIDs); err != nil {
		return fmt.Errorf("seedPostAllowedUsers: %w", err)
	}

	if err := seedPostCategories(db, postIDs); err != nil {
		return fmt.Errorf("seedPostCategories: %w", err)
	}

	commentIDs, err := seedComments(db, userIDs, postIDs)
	if err != nil {
		return fmt.Errorf("seedComments: %w", err)
	}

	if err := seedReactions(db, userIDs, postIDs, commentIDs); err != nil {
		return fmt.Errorf("seedReactions: %w", err)
	}

	groupIDs, err := seedGroups(db, userIDs)
	if err != nil {
		return fmt.Errorf("seedGroups: %w", err)
	}

	if err := seedGroupMembers(db, groupIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupMembers: %w", err)
	}

	if err := seedGroupInvites(db, groupIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupInvites: %w", err)
	}

	if err := seedGroupRequests(db, groupIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupRequests: %w", err)
	}

	if err := seedGroupMessages(db, groupIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupMessages: %w", err)
	}

	eventIDs, err := seedGroupEvents(db, groupIDs, userIDs)
	if err != nil {
		return fmt.Errorf("seedGroupEvents: %w", err)
	}

	if err := seedEventResponses(db, eventIDs, userIDs); err != nil {
		return fmt.Errorf("seedEventResponses: %w", err)
	}

	groupPostIDs, err := seedGroupPosts(db, groupIDs, userIDs)
	if err != nil {
		return fmt.Errorf("seedGroupPosts: %w", err)
	}

	if err := seedGroupPostComments(db, groupPostIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupPostComments: %w", err)
	}

	convIDs, err := seedConversations(db, userIDs)
	if err != nil {
		return fmt.Errorf("seedConversations: %w", err)
	}

	if err := seedMessages(db, convIDs, userIDs); err != nil {
		return fmt.Errorf("seedMessages: %w", err)
	}

	if err := seedNotifications(db, userIDs); err != nil {
		return fmt.Errorf("seedNotifications: %w", err)
	}

	log.Println("[SEED] done.")
	return nil
}

// reset clears seedable tables. CATEGORY is intentionally left alone since
// it's populated by the migration itself (INSERT OR IGNORE).
func reset(db *sql.DB) error {
	fmt.Println("reset db ")
	tables := []string{
		"NOTIFICATIONS", "MESSAGES", "CONVERSATIONS",
		"GROUP_POST_COMMENTS", "GROUP_POSTS",
		"EVENT_RESPONSES", "GROUP_EVENTS", "GROUP_MESSAGES",
		"GROUP_REQUESTS", "GROUP_INVITES", "GROUP_MEMBERS", "GROUPS",
		"COMMENT_REACTIONS", "POST_REACTIONS", "COMMENTS",
		"POST_CATEGORY", "POST_ALLOWED_USERS", "POSTS",
		"FOLLOWS", "SESSIONS", "USERS",
	}
	if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}
	for _, t := range tables {
		if _, err := db.Exec("DELETE FROM " + t); err != nil {
			return fmt.Errorf("delete %s: %w", t, err)
		}
		if _, err := db.Exec("DELETE FROM sqlite_sequence WHERE name = ?", t); err != nil {
			// sqlite_sequence only exists if a table with AUTOINCREMENT was used;
			// ignore the error if it's missing.
			_ = err
		}
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return nil
}

// ---- USERS ----

func seedUsers(db *sql.DB, hashedPW string) ([]int, error) {
	type seedUser struct {
		First, Last, Email, Nickname, About, Birth string
		Avatar                                     string
		Private                                    int
	}
	users := []seedUser{
		{"Alice", "Martin", "alice@example.com", "ali_m", "Coffee & code.", "1996-04-12", "/uploads/seeder/avatars/avatar.jpeg", 0},
		{"Bob", "Nguyen", "bob@example.com", "", "Traveling the world.", "1994-08-23", "/uploads/seeder/avatars/lofi.jpeg", 0},
		{"Chloe", "Dubois", "chloe@example.com", "Chloe", "", "1999-01-05", "/uploads/seeder/avatars/goat.jpg", 1}, // private
		{"David", "Smith", "david@example.com", "david", "Full-stack dev.", "1990-11-30", "", 0},
		{"Emma", "Wilson", "emma@example.com", "", "Photography enthusiast.", "1997-06-18", "", 1}, // private
		{"Farid", "El Amrani", "farid@example.com", "", "Backend > frontend, fight me.", "1993-03-09", "/uploads/seeder/avatars/lofi.jpeg", 0},
		{"Grace", "Lee", "grace@example.com", "", "", "2000-09-27", "", 0},
		{"Hugo", "Costa", "hugo@example.com", "costa77", "Music producer.", "1995-12-14", "", 0},
	}

	ids := make([]int, 0, len(users))
	query := `INSERT INTO USERS (firstname, lastname, email, password, birthdate, nickname, aboutme, avatar, is_private)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, u := range users {
		nickname := sql.NullString{String: u.Nickname, Valid: u.Nickname != ""}
		about := sql.NullString{String: u.About, Valid: u.About != ""}
		avatar := sql.NullString{String: u.Avatar, Valid: u.Avatar != ""}

		res, err := db.Exec(query, u.First, u.Last, u.Email, hashedPW, u.Birth, nickname, about, avatar, u.Private)
		if err != nil {
			return nil, fmt.Errorf("insert user %s: %w", u.Email, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] users: %d\n", len(ids))
	return ids, nil
}

// ---- FOLLOWS ----
// u[0]=Alice u[1]=Bob u[2]=Chloe(private) u[3]=David u[4]=Emma(private) u[5]=Farid u[6]=Grace u[7]=Hugo

func seedFollows(db *sql.DB, u []int) error {
	type f struct {
		Follower, Following int
		Status              string
	}
	follows := []f{
		{u[0], u[1], "accepted"}, // alice -> bob
		{u[1], u[0], "accepted"}, // bob -> alice (mutual)
		{u[0], u[3], "accepted"}, // alice -> david
		{u[3], u[5], "accepted"}, // david -> farid
		{u[5], u[0], "accepted"}, // farid -> alice
		{u[6], u[2], "pending"},  // grace -> chloe (private, awaiting accept)
		{u[7], u[4], "pending"},  // hugo -> emma (private, awaiting accept)
		{u[1], u[4], "accepted"}, // bob -> emma (already accepted earlier)
		{u[2], u[0], "accepted"}, // chloe -> alice
	}
	query := `INSERT INTO FOLLOWS (follower_id, following_id, status) VALUES (?, ?, ?)`
	for _, r := range follows {
		if _, err := db.Exec(query, r.Follower, r.Following, r.Status); err != nil {
			return err
		}
	}
	log.Printf("[SEED] follows: %d\n", len(follows))
	return nil
}

// ---- POSTS ----

func seedPosts(db *sql.DB, u []int) ([]int, error) {
	type p struct {
		UserID      int
		Title, Text string
		Privacy     string
	}
	posts := []p{
		{u[0], "Hello world", "My very first post on this network!", "public"},
		{u[1], "Weekend trip", "Just got back from the mountains 🏔️", "public"},
		{u[3], "", "Debugging is 90% of my job today.", "almost_private"},
		{u[5], "Go tip", "context.Context should be your first param, always.", "public"},
		{u[2], "Private thoughts", "Only a few people should see this.", "private"},
		{u[4], "New camera", "Testing out my new lens today.", "almost_private"},
		{u[7], "New track", "Dropping a new beat this Friday 🎵", "public"},
	}
	query := `INSERT INTO POSTS (user_id, created_at, title, text, image, privacy) VALUES (?, ?, ?, ?, ?, ?)`
	ids := make([]int, 0, len(posts))
	now := time.Now()
	for i, post := range posts {
		title := sql.NullString{String: post.Title, Valid: post.Title != ""}
		createdAt := now.Add(-time.Duration(len(posts)-i) * time.Hour).Format("2006-01-02 15:04:05")
		res, err := db.Exec(query, post.UserID, createdAt, title, post.Text, nil, post.Privacy)
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

// seedPostAllowedUsers grants specific followers access to the 'private' post (posts[4] = Chloe's).
func seedPostAllowedUsers(db *sql.DB, postIDs, u []int) error {
	privatePostID := postIDs[4]  // Chloe's private post
	allowed := []int{u[6], u[0]} // grace, alice
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
	// CATEGORY ids come from the migration's INSERT OR IGNORE; look them up by name.
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

	links := map[int][]int{
		postIDs[0]: {general},
		postIDs[1]: {travel},
		postIDs[3]: {edu},
		postIDs[6]: {ent},
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

// ---- COMMENTS ----

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
		{u[1], postIDs[6], "Can't wait to hear it!"},
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

// ---- REACTIONS ----

func seedReactions(db *sql.DB, u, postIDs, commentIDs []int) error {
	type pr struct {
		UserID, PostID, IsLike int
	}
	postReactions := []pr{
		{u[1], postIDs[0], 1},
		{u[3], postIDs[0], 1},
		{u[5], postIDs[1], 1},
		{u[0], postIDs[3], 1},
		{u[7], postIDs[3], -1},
		{u[0], postIDs[6], 1},
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

// ---- GROUPS ----

func seedGroups(db *sql.DB, u []int) ([]int, error) {
	type g struct {
		CreatorID          int
		Title, Description string
	}
	groups := []g{
		{u[0], "Gophers United", "Everything about Go programming."},
		{u[5], "Weekend Hikers", "Organizing hikes and outdoor trips."},
	}
	query := `INSERT INTO GROUPS (creator_id, title, description) VALUES (?, ?, ?)`
	ids := make([]int, 0, len(groups))
	for _, gr := range groups {
		res, err := db.Exec(query, gr.CreatorID, gr.Title, gr.Description)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] groups: %d\n", len(ids))
	return ids, nil
}

func seedGroupMembers(db *sql.DB, groupIDs, u []int) error {
	type m struct {
		GroupID, UserID int
		Role            string
	}
	members := []m{
		// Gophers United (creator = alice = u[0])
		{groupIDs[0], u[0], "admin"},
		{groupIDs[0], u[3], "member"},
		{groupIDs[0], u[5], "member"},
		{groupIDs[0], u[1], "member"},
		// Weekend Hikers (creator = farid = u[5])
		{groupIDs[1], u[5], "admin"},
		{groupIDs[1], u[0], "member"},
		{groupIDs[1], u[7], "member"},
	}
	query := `INSERT INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, ?)`
	for _, mem := range members {
		if _, err := db.Exec(query, mem.GroupID, mem.UserID, mem.Role); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_members: %d\n", len(members))
	return nil
}

func seedGroupInvites(db *sql.DB, groupIDs, u []int) error {
	type inv struct {
		GroupID, InviterID, InvitedID int
		Status                        string
	}
	invites := []inv{
		{groupIDs[0], u[0], u[6], "pending"}, // alice invites grace to Gophers United
		{groupIDs[1], u[5], u[4], "pending"}, // farid invites emma to Weekend Hikers
	}
	query := `INSERT INTO GROUP_INVITES (group_id, inviter_id, invited_user_id, status) VALUES (?, ?, ?, ?)`
	for _, i := range invites {
		if _, err := db.Exec(query, i.GroupID, i.InviterID, i.InvitedID, i.Status); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_invites: %d\n", len(invites))
	return nil
}

func seedGroupRequests(db *sql.DB, groupIDs, u []int) error {
	type req struct {
		GroupID, UserID int
		Status          string
	}
	requests := []req{
		{groupIDs[0], u[7], "pending"}, // hugo requests to join Gophers United
		{groupIDs[1], u[3], "pending"}, // david requests to join Weekend Hikers
	}
	query := `INSERT INTO GROUP_REQUESTS (group_id, user_id, status) VALUES (?, ?, ?)`
	for _, r := range requests {
		if _, err := db.Exec(query, r.GroupID, r.UserID, r.Status); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_requests: %d\n", len(requests))
	return nil
}

func seedGroupMessages(db *sql.DB, groupIDs, u []int) error {
	type gm struct {
		GroupID, SenderID int
		Text              string
	}
	messages := []gm{
		{groupIDs[0], u[0], "Welcome to Gophers United!"},
		{groupIDs[0], u[3], "Excited to be here 🚀"},
		{groupIDs[0], u[5], "Anyone using generics in prod yet?"},
		{groupIDs[1], u[5], "Hike this Saturday, who's in?"},
		{groupIDs[1], u[0], "Count me in!"},
	}
	query := `INSERT INTO GROUP_MESSAGES (group_id, sender_id, text) VALUES (?, ?, ?)`
	for _, m := range messages {
		if _, err := db.Exec(query, m.GroupID, m.SenderID, m.Text); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_messages: %d\n", len(messages))
	return nil
}

// ---- GROUP EVENTS ----

func seedGroupEvents(db *sql.DB, groupIDs, u []int) ([]int, error) {
	type ev struct {
		GroupID, CreatorID int
		Title, Description string
		EventTime          time.Time
	}
	events := []ev{
		{groupIDs[0], u[0], "Go Meetup #1", "Monthly meetup to discuss Go internals.", time.Now().Add(7 * 24 * time.Hour)},
		{groupIDs[1], u[5], "Sunrise Hike", "Meet at the trailhead at 6am.", time.Now().Add(3 * 24 * time.Hour)},
	}
	query := `INSERT INTO GROUP_EVENTS (group_id, creator_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)`
	ids := make([]int, 0, len(events))
	for _, e := range events {
		res, err := db.Exec(query, e.GroupID, e.CreatorID, e.Title, e.Description, e.EventTime.Format("2006-01-02 15:04:05"))
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] group_events: %d\n", len(ids))
	return ids, nil
}

func seedEventResponses(db *sql.DB, eventIDs, u []int) error {
	type r struct {
		EventID, UserID int
		Status          string
	}
	responses := []r{
		{eventIDs[0], u[0], "going"},
		{eventIDs[0], u[3], "going"},
		{eventIDs[0], u[5], "not_going"},
		{eventIDs[1], u[5], "going"},
		{eventIDs[1], u[0], "going"},
		{eventIDs[1], u[7], "not_going"},
	}
	query := `INSERT INTO EVENT_RESPONSES (event_id, user_id, status) VALUES (?, ?, ?)`
	for _, resp := range responses {
		if _, err := db.Exec(query, resp.EventID, resp.UserID, resp.Status); err != nil {
			return err
		}
	}
	log.Printf("[SEED] event_responses: %d\n", len(responses))
	return nil
}

// ---- GROUP POSTS / COMMENTS (members-only wall) ----

func seedGroupPosts(db *sql.DB, groupIDs, u []int) ([]int, error) {
	type gp struct {
		GroupID, UserID int
		Title, Text     string
	}
	posts := []gp{
		{groupIDs[0], u[0], "Style guide", "Let's agree on gofmt + golangci-lint for the repo."},
		{groupIDs[0], u[5], "", "Anyone tried the new slices package?"},
		{groupIDs[1], u[5], "Trail conditions", "Trail is muddy after yesterday's rain, bring boots."},
	}
	query := `INSERT INTO GROUP_POSTS (group_id, user_id, title, text) VALUES (?, ?, ?, ?)`
	ids := make([]int, 0, len(posts))
	for _, p := range posts {
		title := sql.NullString{String: p.Title, Valid: p.Title != ""}
		res, err := db.Exec(query, p.GroupID, p.UserID, title, p.Text)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] group_posts: %d\n", len(ids))
	return ids, nil
}

func seedGroupPostComments(db *sql.DB, groupPostIDs, u []int) error {
	type c struct {
		GroupPostID, UserID int
		Text                string
	}
	comments := []c{
		{groupPostIDs[0], u[3], "Agreed, let's add it to the CI pipeline."},
		{groupPostIDs[1], u[0], "Yes! Much cleaner than the old sort package."},
		{groupPostIDs[2], u[0], "Thanks for the heads up, bringing boots."},
	}
	query := `INSERT INTO GROUP_POST_COMMENTS (group_post_id, user_id, text) VALUES (?, ?, ?)`
	for _, cm := range comments {
		if _, err := db.Exec(query, cm.GroupPostID, cm.UserID, cm.Text); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_post_comments: %d\n", len(comments))
	return nil
}

// ---- CHAT ----

func seedConversations(db *sql.DB, u []int) ([]int, error) {
	type cv struct {
		User1, User2  int
		LastMessage   string
		LastMessageAt time.Time
	}
	convs := []cv{
		{u[0], u[1], "See you tomorrow!", time.Now().Add(-2 * time.Hour)},
		{u[0], u[5], "Sounds good, thanks!", time.Now().Add(-45 * time.Minute)},
	}
	query := `INSERT INTO CONVERSATIONS (user1_id, user2_id, last_message, last_message_at) VALUES (?, ?, ?, ?)`
	ids := make([]int, 0, len(convs))
	for _, c := range convs {
		res, err := db.Exec(query, c.User1, c.User2, c.LastMessage, c.LastMessageAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] conversations: %d\n", len(ids))
	return ids, nil
}

func seedMessages(db *sql.DB, convIDs, u []int) error {
	type m struct {
		ConversationID, SenderID int
		Text                     string
		IsRead                   int
	}
	messages := []m{
		{convIDs[0], u[0], "Hey Bob, are we still on for tomorrow?", 1},
		{convIDs[0], u[1], "Yep, see you tomorrow!", 0},
		{convIDs[1], u[0], "Can you send me the trail map?", 1},
		{convIDs[1], u[5], "Sounds good, thanks!", 0},
	}
	query := `INSERT INTO MESSAGES (conversation_id, sender_id, text, is_read) VALUES (?, ?, ?, ?)`
	for _, msg := range messages {
		if _, err := db.Exec(query, msg.ConversationID, msg.SenderID, msg.Text, msg.IsRead); err != nil {
			return err
		}
	}
	log.Printf("[SEED] messages: %d\n", len(messages))
	return nil
}

// ---- NOTIFICATIONS ----

func seedNotifications(db *sql.DB, u []int) error {
	type n struct {
		UserID      int
		Type        string
		ReferenceID int
		IsRead      int
	}
	notifications := []n{
		{u[2], "follow_request", u[6], 0},      // chloe: grace requested to follow
		{u[4], "follow_request", u[7], 0},      // emma: hugo requested to follow
		{u[6], "group_invite", u[0], 0},        // grace: invited to Gophers United by alice
		{u[4], "group_invite", u[5], 0},        // emma: invited to Weekend Hikers by farid
		{u[0], "group_join_request", u[7], 0},  // alice (creator): hugo requested to join
		{u[5], "group_join_request", u[3], 0},  // farid (creator): david requested to join
		{u[3], "group_event_created", u[0], 1}, // david: new event in Gophers United (read)
		{u[5], "group_event_created", u[0], 0}, // farid: new event in Gophers United
	}
	query := `INSERT INTO NOTIFICATIONS (user_id, type, reference_id, is_read) VALUES (?, ?, ?, ?)`
	for _, notif := range notifications {
		if _, err := db.Exec(query, notif.UserID, notif.Type, notif.ReferenceID, notif.IsRead); err != nil {
			return err
		}
	}
	log.Printf("[SEED] notifications: %d\n", len(notifications))
	return nil
}
