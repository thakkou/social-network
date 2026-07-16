package repository

import (
	"database/sql"
	"time"
)

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupPost struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
}

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupRepository struct {
	DB *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

func (r *GroupRepository) CreateGroup(g *Group) error {
	query := `INSERT INTO GROUPS (creator_id, title, description) VALUES (?, ?, ?)`
	res, err := r.DB.Exec(query, g.CreatorID, g.Title, g.Description)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	g.ID = int(id)

	_, err = r.DB.Exec(`INSERT INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'admin')`, g.ID, g.CreatorID)
	return err
}

func (r *GroupRepository) RequestToJoin(groupID, userID int) error {
	query := `INSERT INTO GROUP_REQUESTS (group_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := r.DB.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) CreateGroupPost(gp *GroupPost) error {
	query := `INSERT INTO GROUP_POSTS (group_id, user_id, title, text, image) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, gp.GroupID, gp.UserID, gp.Title, gp.Text, gp.Image)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	gp.ID = int(id)
	return nil
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
