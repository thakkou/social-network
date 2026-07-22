package handlers

import (
	"net/http"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
)

type SearchResponse struct {
	Profiles []repository.User  `json:"profiles"`
	Groups   []repository.Group `json:"groups"`
}

// FindQuery searches for users and groups by text query.
// @Summary Search users and groups
// @Description Searches for both users/profiles and groups matching the given text.
// @Tags Search
// @Produce json
// @Param text query string true "Search text"
// @Success 200 {object} SearchResponse "Search results"
// @Router /api/search [get]
func FindQuery(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	userID, _ := middlewares.GetUserID(r)

	profiles, err := Repos.Profile.SearchProfiles(text, userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	groups, err := Repos.Group.SearchGroups(text)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "success", SearchResponse{
		Profiles: profiles,
		Groups:   groups,
	})
}
