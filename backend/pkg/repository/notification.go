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
	query := `UPDATE NOTIFICATIONS SET is_read = 1 WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *NotificationRepository) MarkAllAsReadByUserID(userID int) error {
	query := `UPDATE NOTIFICATIONS SET is_read = 1 WHERE user_id = ? AND is_read = 0`
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *NotificationRepository) DeleteByIDAndUserID(id int, userID int) error {
	query := `DELETE FROM NOTIFICATIONS WHERE id = ? AND user_id = ?`
	res, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	// Optional check: ensure rows were actually modified (proves ownership)
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return errors.New("notification not found or unauthorized")
	}

	return nil
}

func (r *NotificationRepository) DeleteAllByUserID(userID int) error {
	query := `DELETE FROM NOTIFICATIONS WHERE user_id = ?`
	_, err := r.DB.Exec(query, userID)
	return err
}
