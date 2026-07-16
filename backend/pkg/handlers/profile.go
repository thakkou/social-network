package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	repModal "01social/pkg/models/representation"
	"01social/pkg/repository"
	"01social/pkg/representation"
	"01social/pkg/utilities"
)

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

	followStatus, err := followRepo.GetFollowStatus(userID, profileId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var user *repository.User
	var profileRes repModal.ProfileResponse

	if followStatus == "accepted" {
		fmt.Println("user following the profile")
	}
	if !private || followStatus == "accepted" {

		user, err = profileRepo.GetProfile(profileId, false)

		following, _ := followRepo.GetFollowing(profileId)
		fmt.Println("user ", profileId, "follow", following, len(following))
		followers, _ := followRepo.GetFollowers(profileId)
		fmt.Println("user ", profileId, "followers", followers, len(followers))

		fmt.Println("have right to get the data")

		profileRes = representation.UserToProfileResponse(user, followers, following, followStatus)

	} else {
		user, err = profileRepo.GetProfile(profileId, true)
		profileRes = representation.UserToProfileResponse(user, nil, nil, followStatus)
		fmt.Println("get only public data")
	}

	utilities.WriteJSON(w, http.StatusForbidden, "IDs are different", profileRes)
}
