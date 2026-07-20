package repository

import (
	"time"

	"01social/pkg/utilities"
)

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *GroupRepository) ListGroupEvents(groupID int) ([]GroupEvent, error) {
	rows, err := r.DB.Query(`SELECT id, group_id, creator_id, title, description, event_time, created_at FROM GROUP_EVENTS WHERE group_id = ? ORDER BY event_time ASC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]GroupEvent, 0)
	for rows.Next() {
		var e GroupEvent
		var eventTime string
		var createdAt string
		if err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &eventTime, &createdAt); err != nil {
			return nil, err
		}
		e.EventTime, _ = time.Parse("2006-01-02 15:04:05", eventTime)
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *GroupRepository) CreateEvent(e *GroupEvent) error {
	query := `INSERT INTO GROUP_EVENTS (group_id, creator_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, e.GroupID, e.CreatorID, e.Title, e.Description, e.EventTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	e.ID = int(id)
	return nil
}

func (r *GroupRepository) RespondToEvent(eventID, userID int, status string) error {
	query := `INSERT INTO EVENT_RESPONSES (event_id, user_id, status) VALUES (?, ?, ?)
	          ON CONFLICT(event_id, user_id) DO UPDATE SET status = ?`
	_, err := r.DB.Exec(query, eventID, userID, status, status)
	return err
}

// EventResponder is one user's response ("going", "not_going", etc.) to an event,
// enriched with the user's profile details.
type EventResponder struct {
	UserID    int    `json:"user_id"`
	Status    string `json:"status"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Avatar    string `json:"avatar"`
}

// GetEventsResponses batch-fetches who responded to each event and how,
// for a set of event IDs. Returns a map keyed by event ID; events with no
// responses yet are simply absent from the map.
func (r *GroupRepository) GetEventsResponses(eventIDs []int) (map[int][]EventResponder, error) {
	responses := make(map[int][]EventResponder, len(eventIDs))
	if len(eventIDs) == 0 {
		return responses, nil
	}

	placeholders, args := utilities.PlaceholdersForInts(eventIDs)

	query := `
SELECT
    er.event_id,
    er.user_id,
    er.status,
    COALESCE(u.nickname, '') AS nickname,
    u.firstname,
    u.lastname,
    COALESCE(u.avatar, '') AS avatar
FROM EVENT_RESPONSES er
JOIN USERS u ON u.id = er.user_id
WHERE er.event_id IN (` + placeholders + `)
ORDER BY er.event_id
`
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventID int
		var responder EventResponder
		if err := rows.Scan(
			&eventID,
			&responder.UserID,
			&responder.Status,
			&responder.Nickname,
			&responder.Firstname,
			&responder.Lastname,
			&responder.Avatar,
		); err != nil {
			return nil, err
		}
		responses[eventID] = append(responses[eventID], responder)
	}
	return responses, rows.Err()
}
