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
	return ticket, nil
}

func RedeemTicket(ticket string) (int, error) {
	var userID int
	cutoff := time.Now().Add(-ticketTTL)

	err := db.Database.QueryRow(
		"SELECT user_id FROM WS_TICKETS WHERE ticket = ? AND created_at > ?",
		ticket, cutoff,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("invalid or expired ticket: %w", err)
	}

	_, _ = db.Database.Exec("DELETE FROM WS_TICKETS WHERE ticket = ?", ticket)

	return userID, nil
}
