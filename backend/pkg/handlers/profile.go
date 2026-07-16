package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
)

/*
data gets
-see first and last name,is_private
followed him or public
--see following,followers posts
*/
func GetProfile(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusNotFound, "404 not found", nil)
		return
	}

	// Check path: /api/profile/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "api" || parts[1] != "profile" {
		utilities.WriteJSON(w, http.StatusNotFound, "404 not found", nil)
		return
	}

	profileId, err := strconv.Atoi(parts[2])
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	// Get authenticated user ID
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	// Compare IDs
	if profileId == userID {
		utilities.WriteJSON(w, http.StatusOK, "you are the same", nil)
		return
	}
	// is the user private
	userRepo := repository.NewUserRepository(db.Database)
	private, err := userRepo.IsPrivateUser(profileId)
	if err != nil {
		fmt.Println("User not found:", err)
		return
	}

	if private {
		fmt.Println("This user is private")
	} else {
		fmt.Println("This user is public")
	}
	followRepo := repository.NewFollowRepository(db.Database)
	profileRepo := repository.NewProfileRepository(db.Database)

	isFlollowing, err := followRepo.IsFollowing(userID, profileId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var user *repository.User

	if isFlollowing {
		fmt.Println("user following the profile")
	}
	if !private || isFlollowing {

		user, err = profileRepo.GetProfile(profileId, false)

		fmt.Println("have right to get the data")
	} else {
		user, err = profileRepo.GetProfile(profileId, true)

		fmt.Println("get only public data")
	}

	utilities.WriteJSON(w, http.StatusForbidden, "IDs are different", user)
}
