package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/representation"
	"01social/pkg/utilities"
)

func UpdateProfilePrivacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	var payload struct {
		IsPrivate int `json:"is_private"`
	}
	if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if payload.IsPrivate != 0 && payload.IsPrivate != 1 {
		utilities.WriteJSON(w, http.StatusBadRequest, "is_private must be 0 or 1", nil)
		return
	}

	if err := Repos.User.UpdatePrivacy(userID, payload.IsPrivate); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not update profile privacy", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "profile privacy updated", map[string]any{"is_private": payload.IsPrivate})
}

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
	if userID == profileID {
		user, err := Repos.Profile.GetProfile(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusNotFound, "profile not found", nil)
			return
		}
		followers, err := Repos.Follow.GetFollowers(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get followers", nil)
			return
		}
		following, err := Repos.Follow.GetFollowing(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get following", nil)
			return
		}
		posts, err := Repos.Post.GetPostsUserID(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get posts", nil)
			return
		}
		profileRes := representation.UserToProfileResponse(user, followers, following, posts, "accepted")
		utilities.WriteJSON(w, http.StatusOK, "profile data", profileRes)
		return
	}

	// Check if the requested profile belongs to a private account
	private, err := Repos.User.IsPrivateUser(profileID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "user not found", nil)
		return
	}

	// Get the relationship status between current user and profile owner
	// Possible values: accepted, pending, none
	followStatus, err := Repos.Follow.GetFollowStatus(userID, profileID)
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
		user, err = Repos.Profile.GetProfile(profileID)
	} else {
		user, err = Repos.Profile.GetPublicProfile(profileID)
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
		followers, err = Repos.Follow.GetFollowers(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get followers", nil)
			return
		}

		// Get users this profile follows
		following, err = Repos.Follow.GetFollowing(profileID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get following", nil)
			return
		}
		posts, err = Repos.Post.GetPostsUserID(profileID)
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
