package utilities

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ── In-memory ticket store ──

type wsTicket struct {
	userID  int
	expires time.Time
}

var (
	ticketMu    sync.RWMutex
	tickets     = make(map[string]wsTicket)
	ticketTTL   = 30 * time.Second
	cleanupOnce sync.Once
)

func startTicketCleanup() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ticketMu.Lock()
			now := time.Now()
			for key, t := range tickets {
				if now.After(t.expires) {
					delete(tickets, key)
				}
			}
			ticketMu.Unlock()
		}
	}()
}

func CreateTicket(userID int) (string, error) {
	cleanupOnce.Do(startTicketCleanup)

	ticket := uuid.New().String()
	ticketMu.Lock()
	tickets[ticket] = wsTicket{
		userID:  userID,
		expires: time.Now().Add(ticketTTL),
	}
	ticketMu.Unlock()

	return ticket, nil
}

func RedeemTicket(ticket string) (int, error) {
	ticketMu.Lock()
	defer ticketMu.Unlock()

	t, ok := tickets[ticket]
	if !ok {
		return 0, fmt.Errorf("ticket not found")
	}
	if time.Now().After(t.expires) {
		delete(tickets, ticket)
		return 0, fmt.Errorf("ticket expired")
	}

	delete(tickets, ticket)
	return t.userID, nil
}
