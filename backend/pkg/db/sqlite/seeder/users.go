package seeder

import (
	"database/sql"
	"fmt"
	"log"
)

type seedUser struct {
	First, Last, Email, Nickname, About, Birth string
	Avatar                                     string
	Private                                    int
}

var seedUsersData = []seedUser{
	{"Alice", "Martin", "alice@example.com", "ali_m", "Coffee & code.", "1996-04-12", "/uploads/seeder/avatars/avatar.jpeg", 0},
	{"Bob", "Nguyen", "bob@example.com", "", "Traveling the world.", "1994-08-23", "/uploads/seeder/avatars/lofi.jpeg", 0},
	{"Chloe", "Dubois", "chloe@example.com", "Chloe", "", "1999-01-05", "/uploads/seeder/avatars/goat.jpg", 1},
	{"David", "Smith", "david@example.com", "david", "Full-stack dev.", "1990-11-30", "", 0},
	{"Emma", "Wilson", "emma@example.com", "", "Photography enthusiast.", "1997-06-18", "", 1},
	{"Farid", "El Amrani", "farid@example.com", "", "Backend > frontend, fight me.", "1993-03-09", "/uploads/seeder/avatars/lofi.jpeg", 0},
	{"Grace", "Lee", "grace@example.com", "", "", "2000-09-27", "", 0},
	{"Hugo", "Costa", "hugo@example.com", "costa77", "Music producer.", "1995-12-14", "", 0},
}

func seedUsers(db *sql.DB, hashedPW string) ([]int, error) {
	ids := make([]int, 0, len(seedUsersData))
	query := `INSERT INTO USERS (firstname, lastname, email, password, birthdate, nickname, aboutme, avatar, is_private) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, u := range seedUsersData {
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

		// Additional follows for richer data
		{u[0], u[6], "accepted"}, // alice -> grace
		{u[1], u[3], "accepted"}, // bob -> david
		{u[7], u[0], "accepted"}, // hugo -> alice
		{u[3], u[0], "accepted"}, // david -> alice
		{u[5], u[7], "accepted"}, // farid -> hugo
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
