package handlers

import (
	"net/http"
	"strconv"

	"01social/pkg/middlewares"
	"01social/pkg/utilities"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	notes, err := Repos.Notification.GetByUserID(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch notifications", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notifications fetched", notes)
}

func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid notification id", nil)
		return
	}

	if err := Repos.Notification.MarkAsRead(id); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not mark notification as read", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notification marked as read", nil)
}
