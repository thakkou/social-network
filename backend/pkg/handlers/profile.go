package handlers

import (
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/representation"
	"01social/pkg/utilities"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests for fetching profiles
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusNotFound, "404 not found", nil)
		return
	}

	// Validate URL format: /api/profile/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "api" || parts[1] != "profile" {
		utilities.WriteJSON(w, http.StatusNotFound, "404 not found", nil)
		return
	}

	// Convert profile ID from URL string to integer
	profileID, err := strconv.Atoi(parts[2])
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	// Get the currently authenticated user's ID
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Database)
	followRepo := repository.NewFollowRepository(db.Database)
	profileRepo := repository.NewProfileRepository(db.Database)
	postRepo := repository.NewPostRepository(db.Database)

	// Check if the requested profile belongs to a private account
	private, err := userRepo.IsPrivateUser(profileID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "user not found", nil)
		return
	}

	// Get the relationship status between current user and profile owner
	// Possible values: accepted, pending, none
	followStatus, err := followRepo.GetFollowStatus(userID, profileID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not get follow status", nil)
		return
	}

	// User can access the complete profile if:
	// - The account is public
	// - The current user is an accepted follower
	canSeeFullProfile := !private || followStatus == "accepted"

	var user *repository.User

	// Fetch complete profile data or only public information
	if canSeeFullProfile {
		user, err = profileRepo.GetProfile(profileID, false)
	} else {
		user, err = profileRepo.GetProfile(profileID, true)
	}

	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "profile not found", nil)
		return
	}

	// Prepare followers and following lists
	// They are only visible when the user has permission to see the full profile
	var followers []repository.User
	var following []repository.User
	var posts []repository.Post

	if canSeeFullProfile {

		// Get users who follow this profile
		followers, err = followRepo.GetFollowers(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get followers", nil)
			return
		}

		// Get users this profile follows
		following, err = followRepo.GetFollowing(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get following", nil)
			return
		}
		// get posts for this profile
		// get posts for this profile
		posts, err = postRepo.GetPostsUserID(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get posts", nil)
			return
		}
	}

	// Convert database models into API response model
	profileRes := representation.UserToProfileResponse(
		user,
		followers,
		following,
		posts,
		followStatus,
	)

	// Return profile data
	utilities.WriteJSON(w, http.StatusOK, "profile data", profileRes)
}
