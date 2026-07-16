package handlers

import (
	"net/http"

	"01social/pkg/utilities"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, 200, "user data  get succes", nil)
}
