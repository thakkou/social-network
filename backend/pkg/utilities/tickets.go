package utilities

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	db "01social/pkg/db/sqlite"
)

const ticketTTL = 30 * time.Second

func CreateTicket(userID int) (string, error) {
	ticket := uuid.New().String()
	_, err := db.Database.Exec(
		"INSERT INTO WS_TICKETS (ticket, user_id) VALUES (?, ?)",
		ticket, userID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create ws ticket: %w", err)
	}

	fmt.Println("Inserted ticket:", ticket)

	rows, _ := db.Database.Query("SELECT ticket, user_id, created_at FROM WS_TICKETS")
	defer rows.Close()

	for rows.Next() {
		var t string
		var id int
		var created string
		rows.Scan(&t, &id, &created)
		fmt.Println("DB:", t, id, created)
	}
	return ticket, nil
}

func RedeemTicket(ticket string) (int, error) {
	var userID int

	err := db.Database.QueryRow(
		"SELECT user_id FROM WS_TICKETS WHERE ticket = ?",
		ticket,
	).Scan(&userID)
	if err != nil {
		fmt.Println("RedeemTicket:", err)
		return 0, err
	}

	_, _ = db.Database.Exec(
		"DELETE FROM WS_TICKETS WHERE ticket = ?",
		ticket,
	)

	return userID, nil
}
