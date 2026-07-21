package handlers

import (
	"net/http"

	"01social/pkg/middlewares"
	"01social/pkg/utilities"
)

func CreateWsTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not authenticated", nil)
		return
	}

	ticket, err := utilities.CreateTicket(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to create ticket", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "ticket created", map[string]string{
		"ticket": ticket,
	})
}
