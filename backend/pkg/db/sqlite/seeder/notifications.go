package seeder

import (
	"database/sql"
	"log"
)

func seedNotifications(db *sql.DB, u []int) error {
	type n struct {
		UserID     int
		ActorID    int
		Type       string
		ObjectType string
		ObjectID   int
		IsRead     int
	}
	notifications := []n{
		{UserID: u[2], ActorID: u[6], Type: "follow_request", ObjectType: "follow", ObjectID: u[6], IsRead: 0},
		{UserID: u[0], ActorID: u[1], Type: "post_reaction", ObjectType: "post", ObjectID: 1, IsRead: 0},
		{UserID: u[0], ActorID: u[3], Type: "comment", ObjectType: "comment", ObjectID: 2, IsRead: 0},
		{UserID: u[0], ActorID: u[2], Type: "follow_accepted", ObjectType: "follow", ObjectID: u[2], IsRead: 1},
		{UserID: u[0], ActorID: u[5], Type: "group_invite", ObjectType: "group_invite", ObjectID: 2, IsRead: 1},
		{UserID: u[1], ActorID: u[0], Type: "post_reaction", ObjectType: "post", ObjectID: 2, IsRead: 0},
		{UserID: u[1], ActorID: u[2], Type: "comment", ObjectType: "comment", ObjectID: 1, IsRead: 1},
		{UserID: u[1], ActorID: u[6], Type: "follow_request", ObjectType: "follow", ObjectID: u[6], IsRead: 0},
		{UserID: u[1], ActorID: u[0], Type: "message", ObjectType: "conversation", ObjectID: 1, IsRead: 0},
		{UserID: u[2], ActorID: u[1], Type: "post_reaction", ObjectType: "post", ObjectID: 9, IsRead: 0},
		{UserID: u[2], ActorID: u[0], Type: "comment", ObjectType: "comment", ObjectID: 5, IsRead: 0},
		{UserID: u[2], ActorID: u[0], Type: "follow_accepted", ObjectType: "follow", ObjectID: u[0], IsRead: 1},
		{UserID: u[2], ActorID: u[5], Type: "group_event", ObjectType: "event", ObjectID: 2, IsRead: 0},
		{UserID: u[0], ActorID: u[7], Type: "follow_request", ObjectType: "follow", ObjectID: u[7], IsRead: 0},
		{UserID: u[1], ActorID: u[0], Type: "group_event", ObjectType: "event", ObjectID: 1, IsRead: 0},
		{UserID: u[4], ActorID: u[7], Type: "follow_request", ObjectType: "follow", ObjectID: u[7], IsRead: 0},
		{UserID: u[6], ActorID: u[0], Type: "group_invite", ObjectType: "group_invite", ObjectID: 1, IsRead: 0},
		{UserID: u[4], ActorID: u[5], Type: "group_invite", ObjectType: "group_invite", ObjectID: 2, IsRead: 0},
		{UserID: u[0], ActorID: u[7], Type: "group_join_request", ObjectType: "group_request", ObjectID: 1, IsRead: 0},
		{UserID: u[5], ActorID: u[3], Type: "group_join_request", ObjectType: "group_request", ObjectID: 2, IsRead: 0},
		{UserID: u[3], ActorID: u[0], Type: "group_event", ObjectType: "event", ObjectID: 1, IsRead: 1},
		{UserID: u[5], ActorID: u[0], Type: "group_event", ObjectType: "event", ObjectID: 1, IsRead: 0},
	}
	query := `INSERT INTO NOTIFICATIONS (user_id, actor_id, type, object_type, object_id, is_read) VALUES (?, ?, ?, ?, ?, ?)`
	for _, notif := range notifications {
		if _, err := db.Exec(query, notif.UserID, notif.ActorID, notif.Type, notif.ObjectType, notif.ObjectID, notif.IsRead); err != nil {
			return err
		}
	}
	log.Printf("[SEED] notifications: %d\n", len(notifications))
	return nil
}
