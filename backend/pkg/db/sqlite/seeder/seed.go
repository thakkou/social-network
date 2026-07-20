package seeder

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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

	if err := seedGroupPostReactions(db, groupPostIDs, userIDs); err != nil {
		return fmt.Errorf("seedGroupPostReactions: %w", err)
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

func reset(db *sql.DB) error {
	fmt.Println("reset db ")
	tables := []string{
		"NOTIFICATIONS", "MESSAGES", "CONVERSATIONS",
		"GROUP_POST_COMMENTS", "GROUP_POSTS",
		"EVENT_RESPONSES", "GROUP_EVENTS", "GROUP_MESSAGES",
		"GROUP_REQUESTS", "GROUP_INVITES", "GROUP_MEMBERS", "GROUPS",
		"COMMENT_REACTIONS", "POST_REACTIONS", "COMMENTS",
		"POST_CATEGORY", "POST_ALLOWED_USERS", "POSTS", "GROUP_POST_REACTIONS",
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
			_ = err
		}
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return nil
}
