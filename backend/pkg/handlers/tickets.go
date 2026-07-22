package handlers

import (
	"fmt"
	"net/http"

	"01social/pkg/middlewares"
	"01social/pkg/utilities"
)

// CreateWsTicket creates a one-time use ticket for WebSocket authentication.
// @Summary Create WebSocket ticket
// @Description Creates a one-time use ticket that can be exchanged for a WebSocket connection. Tickets expire after a short time.
// @Tags WebSocket
// @Produce json
// @Success 200 {object} map[string]string "Ticket created"
// @Failure 401 {object} map[string]string "Not authenticated"
// @Router /api/ws-ticket [post]
func CreateWsTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("creat a ws ticke")
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		fmt.Println("not ok")
		utilities.WriteJSON(w, http.StatusUnauthorized, "not authenticated", nil)
		return
	}

	ticket, err := utilities.CreateTicket(userID)
	if err != nil {
		fmt.Println("enable creatng aws")
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to create ticket", nil)
		return
	}
	fmt.Println("done creating the tickets")
	utilities.WriteJSON(w, http.StatusOK, "ticket created", map[string]string{
		"ticket": ticket,
	})
}
