package handlers

import (
	"net/http"

	"01social/pkg/middlewares"
	"01social/pkg/utilities"
)

/*
data gets
-see first and last name,is_private
followed him or public
--see following,followers posts
*/
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middlewares.GetUserID(r)
	utilities.WriteJSON(w, 200, "user data  get succes", userID)
}
