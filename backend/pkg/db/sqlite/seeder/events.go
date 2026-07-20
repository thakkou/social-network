package seeder

import (
	"database/sql"
	"log"
	"time"
)

func seedGroupEvents(db *sql.DB, groupIDs, u []int) ([]int, error) {
	type ev struct {
		GroupID, CreatorID int
		Title, Description string
		EventTime          time.Time
	}
	events := []ev{
		{groupIDs[0], u[0], "Go Meetup", "Monthly Go developers meetup.", time.Now().Add(7 * 24 * time.Hour)},
		{groupIDs[1], u[5], "Football Match", "Friendly football game.", time.Now().Add(3 * 24 * time.Hour)},
		{groupIDs[2], u[7], "Museum Visit", "Visit the modern art museum together.", time.Now().Add(10 * 24 * time.Hour)},
		{groupIDs[3], u[1], "Weekend Road Trip", "Two-day trip to the mountains.", time.Now().Add(14 * 24 * time.Hour)},
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
		{eventIDs[0], u[0], "going"}, {eventIDs[0], u[3], "going"}, {eventIDs[0], u[5], "not_going"},
		{eventIDs[1], u[5], "going"}, {eventIDs[1], u[0], "going"}, {eventIDs[1], u[7], "not_going"},
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
