package repository

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Type        string    `json:"type"`
	ReferenceID int       `json:"reference_id"`
	IsRead      int       `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message,omitempty"`
	Title       string    `json:"title,omitempty"`
}

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(n *Notification) error {
	query := `INSERT INTO NOTIFICATIONS (user_id, type, reference_id) VALUES (?, ?, ?)`
	res, err := r.DB.Exec(query, n.UserID, n.Type, n.ReferenceID)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	n.ID = int(id)
	return nil
}

func (r *NotificationRepository) GetByUserID(userID int) ([]Notification, error) {
	query := `SELECT id, user_id, type, reference_id, is_read, created_at FROM NOTIFICATIONS WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Notification
	for rows.Next() {
		var n Notification
		var tStr string
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.ReferenceID, &n.IsRead, &tStr); err != nil {
			return nil, err
		}
		n.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", tStr)
		n.Message, n.Title = r.describeNotification(n.Type, n.ReferenceID)
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *NotificationRepository) describeNotification(typ string, referenceID int) (string, string) {
	switch typ {
	case "follow_request":
		return "Someone requested to follow you.", "Follow request"
	case "group_invite":
		return "You were invited to a group.", "Group invite"
	case "group_join_request":
		return "A user requested to join your group.", "Group join request"
	case "group_event":
		return "A new event was created in a group you belong to.", "Group event"
	case "group_message":
		return "A new message was posted in a group chat.", "Group message"
	default:
		return "You have a new notification.", "Notification"
	}
}

func (r *NotificationRepository) MarkAsRead(id int) error {
	_, err := r.DB.Exec("UPDATE NOTIFICATIONS SET is_read = 1 WHERE id = ?", id)
	return err
}
