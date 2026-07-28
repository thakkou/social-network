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

// UpdateProfilePrivacy toggles the profile privacy setting.
// @Summary Update profile privacy
// @Description Sets the profile to private (1) or public (0). Private profiles require follow requests.
// @Tags Profile
// @Accept json
// @Produce json
// @Param payload body object true "Privacy setting" SchemaExample({"is_private":1})
// @Success 200 {object} map[string]any "Profile privacy updated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/profile/privacy [put]
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

// GetProfile fetches a user's profile by ID.
// @Summary Get user profile
// @Description Returns profile data for a user. For private profiles, only public info is shown unless the current user is a follower.
// @Tags Profile
// @Produce json
// @Param id path int true "Profile/User ID"
// @Success 200 {object} repModal.ProfileResponse "Profile data fetched"
// @Failure 401 {object} map[string]string "Not logged in"
// @Failure 404 {object} map[string]string "Profile not found"
// @Router /api/profile/{id} [get]
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

// UpdateProfile updates the current user's profile (nickname, about me, avatar).
// @Summary Update own profile
// @Description Updates the logged-in user's nickname, about me section, and optional avatar image.
// @Tags Profile
// @Accept mpfd
// @Produce json
// @Param nickname formData string false "New nickname"
// @Param aboutme formData string false "New about me text"
// @Param avatar formData file false "New avatar image"
// @Success 200 {object} map[string]string "Profile updated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/profile/update [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	var nickname, aboutme string
	var avatarPath string

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		var payload struct {
			Nickname string `json:"nickname"`
			AboutMe  string `json:"aboutme"`
		}
		if jsonErr := utilities.ReadJSONRequestIntoStruct(r, &payload); jsonErr != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		nickname = strings.TrimSpace(payload.Nickname)
		aboutme = strings.TrimSpace(payload.AboutMe)
	} else {
		nickname = strings.TrimSpace(r.FormValue("nickname"))
		aboutme = strings.TrimSpace(r.FormValue("aboutme"))

		if file, header, err := r.FormFile("avatar"); err == nil {
			defer file.Close()
			if saved, saveErr := utilities.SaveImage(file, header, "uploads/avatars/"); saveErr == nil {
				avatarPath = saved
			}
		}
	}

	// ── Validate nickname ──
	if nickname != "" && !utilities.IsValidName(nickname) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid nickname: use 2–50 characters, letters, numbers, underscores, hyphens, apostrophes, and periods only", nil)
		return
	}

	// ── Validate about me ──
	if aboutme != "" && !utilities.IsValidDescription(aboutme) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid 'about me': must be 2048 characters or less", nil)
		return
	}

	if err := Repos.User.UpdateProfile(userID, nickname, aboutme, avatarPath); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not update profile", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "profile updated", nil)
}
