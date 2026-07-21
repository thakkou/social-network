package handlers

import (
	"fmt"
	"net/http"

	"01social/pkg/middlewares"
	"01social/pkg/utilities"
)

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
