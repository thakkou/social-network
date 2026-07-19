package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Notification struct {
	ID         int       `json:"id"`
	UserID     int       `json:"id"`
	ActorID    int       `json:"actor_id"`
	Type       string    `json:"type"`
	ObjectType string    `json:"object_type"`
	ObjectID   int       `json:"object_id"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(n *Notification) error {
	query := `
	INSERT INTO NOTIFICATIONS
	(user_id, actor_id, type, object_type, object_id)
	VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.DB.Exec(
		query,
		n.UserID,
		n.ActorID,
		n.Type,
		n.ObjectType,
		n.ObjectID,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	n.ID = int(id)

	return nil
}

func (r *NotificationRepository) GetByUserID(userID int, typeNotif string) ([]Notification, error) {
	if typeNotif != "unread" && typeNotif != "all" {
		return nil, errors.New("invalid notification type")
	}
	fmt.Println("get notifs", typeNotif)

	query := `
	SELECT
		id,
		actor_id,
		type,
		object_type,
		object_id,
		is_read,
		created_at
	FROM NOTIFICATIONS
	WHERE user_id = ?
	`

	args := []interface{}{userID}

	if typeNotif == "unread" {
		query += " AND is_read = 0"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]Notification, 0)

	for rows.Next() {
		var n Notification

		err := rows.Scan(
			&n.ID,
			&n.ActorID,
			&n.Type,
			&n.ObjectType,
			&n.ObjectID,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkAsRead(id int) error {
	_, err := r.DB.Exec(
		"UPDATE NOTIFICATIONS SET is_read = 1 WHERE id = ?",
		id,
	)

	return err
}
