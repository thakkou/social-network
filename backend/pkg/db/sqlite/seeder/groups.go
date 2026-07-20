package seeder

import (
	"database/sql"
	"log"
)

func seedGroups(db *sql.DB, u []int) ([]int, error) {
	type g struct {
		CreatorID   int
		Title       string
		Description string
		Logo        string
		Background  string
	}
	groups := []g{
		{CreatorID: u[0], Title: "Gophers United", Description: "Everything about Go programming, backend development, concurrency and open source.", Logo: "/uploads/seeder/groups/golage-log.jpeg", Background: "/uploads/seeder/groups/golang-bg.jpeg"},
		{CreatorID: u[5], Title: "Sports Club", Description: "Football, basketball, running, fitness and every kind of sport.", Logo: "/uploads/seeder/groups/sport-logo.jpeg", Background: "/uploads/seeder/groups/sport-bg.jpeg"},
		{CreatorID: u[7], Title: "Culture & Arts", Description: "Books, music, cinema, painting and cultural events.", Logo: "/uploads/seeder/groups/cultur-log.jpeg", Background: "/uploads/seeder/groups/culture-bg.jpeg"},
		{CreatorID: u[1], Title: "Travel Explorers", Description: "Share destinations, travel tips, hiking adventures and unforgettable experiences.", Logo: "/uploads/seeder/groups/travel-logo.jpeg", Background: "/uploads/seeder/groups/travel-bg.jpeg"},
	}
	query := `INSERT INTO GROUPS (creator_id, title, description, logo, background) VALUES (?, ?, ?, ?, ?)`
	ids := make([]int, 0, len(groups))
	for _, gr := range groups {
		res, err := db.Exec(query, gr.CreatorID, gr.Title, gr.Description, gr.Logo, gr.Background)
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
		{groupIDs[0], u[0], "admin"}, {groupIDs[0], u[1], "member"}, {groupIDs[0], u[3], "member"}, {groupIDs[0], u[5], "member"},
		{groupIDs[1], u[5], "admin"}, {groupIDs[1], u[0], "member"}, {groupIDs[1], u[6], "member"}, {groupIDs[1], u[7], "member"},
		{groupIDs[2], u[7], "admin"}, {groupIDs[2], u[2], "member"}, {groupIDs[2], u[4], "member"}, {groupIDs[2], u[1], "member"},
		{groupIDs[3], u[1], "admin"}, {groupIDs[3], u[0], "member"}, {groupIDs[3], u[4], "member"}, {groupIDs[3], u[6], "member"},
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
		{groupIDs[0], u[0], u[6], "pending"},
		{groupIDs[1], u[5], u[4], "pending"},
		{groupIDs[2], u[7], u[3], "pending"},
		{groupIDs[3], u[1], u[2], "pending"},
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
		{groupIDs[0], u[7], "pending"},
		{groupIDs[1], u[3], "pending"},
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
		{groupIDs[0], u[5], "Anyone using Go 1.25?"},
		{groupIDs[0], u[3], "Concurrency is amazing."},
		{groupIDs[1], u[5], "Football match this Friday?"},
		{groupIDs[1], u[7], "I'm in! ⚽"},
		{groupIDs[1], u[0], "See you at 7 PM."},
		{groupIDs[2], u[7], "Movie night this weekend?"},
		{groupIDs[2], u[1], "Interstellar gets my vote."},
		{groupIDs[2], u[2], "I'd rather visit a museum."},
		{groupIDs[3], u[1], "Best destination for summer?"},
		{groupIDs[3], u[6], "I recommend Morocco 🇲🇦"},
		{groupIDs[3], u[4], "Japan is on my bucket list."},
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

func seedGroupPosts(db *sql.DB, groupIDs, u []int) ([]int, error) {
	type gp struct {
		GroupID, UserID int
		Title, Text     string
	}
	posts := []gp{
		{groupIDs[0], u[0], "Favorite Go Feature", "Mine is goroutines."},
		{groupIDs[0], u[5], "", "Who's using generics?"},
		{groupIDs[1], u[5], "Weekend Match", "Who's available this Saturday?"},
		{groupIDs[1], u[7], "", "I'll reserve the field."},
		{groupIDs[2], u[7], "Best Movie", "Recommend your favorite movie."},
		{groupIDs[2], u[1], "", "I'm reading Dune right now."},
		{groupIDs[3], u[1], "Dream Destination", "Where do you want to travel next?"},
		{groupIDs[3], u[4], "", "I want to visit Iceland."},
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
		{groupPostIDs[0], u[3], "Goroutines changed the way I write code."},
		{groupPostIDs[1], u[0], "Generics are finally mature."},
		{groupPostIDs[2], u[6], "Count me in!"},
		{groupPostIDs[3], u[5], "Perfect."},
		{groupPostIDs[4], u[2], "The Godfather never gets old."},
		{groupPostIDs[5], u[7], "Great choice!"},
		{groupPostIDs[6], u[6], "Japan is also my dream destination."},
		{groupPostIDs[7], u[1], "Iceland looks incredible in winter."},
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

func seedGroupPostReactions(db *sql.DB, groupPostIDs, u []int) error {
	type reaction struct {
		GroupPostID int
		UserID      int
		IsLike      int
	}
	reactions := []reaction{
		{groupPostIDs[0], u[1], 1}, {groupPostIDs[0], u[3], 1}, {groupPostIDs[0], u[5], -1},
		{groupPostIDs[1], u[0], 1}, {groupPostIDs[1], u[3], 1},
		{groupPostIDs[2], u[0], 1}, {groupPostIDs[2], u[7], 1},
		{groupPostIDs[3], u[5], 1},
		{groupPostIDs[4], u[2], 1},
		{groupPostIDs[5], u[7], -1},
		{groupPostIDs[6], u[6], 1},
		{groupPostIDs[7], u[1], 1},
	}
	query := `INSERT INTO GROUP_POST_REACTIONS (group_post_id, user_id, is_like) VALUES (?, ?, ?)`
	for _, r := range reactions {
		if _, err := db.Exec(query, r.GroupPostID, r.UserID, r.IsLike); err != nil {
			return err
		}
	}
	log.Printf("[SEED] group_post_reactions: %d\n", len(reactions))
	return nil
}


