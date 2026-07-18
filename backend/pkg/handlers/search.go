package handlers

import (
	"net/http"

	"01social/pkg/repository"
	"01social/pkg/utilities"
)

type SearchResponse struct {
	Profiles []repository.User  `json:"profiles"`
	Groups   []repository.Group `json:"groups"`
}

func FindQuery(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	profiles, err := Repos.Profile.SearchProfiles(text)
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
