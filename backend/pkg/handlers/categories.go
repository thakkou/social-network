package handlers

import (
	"fmt"
	"net/http"

	"01social/pkg/db/sqlite"
	"01social/pkg/utilities"
)

func GetCategoriesByPost(postId int) ([]string, error) {
	var categories []string

	rows, err := sqlite.DB().Query(`
		SELECT c.name
		FROM category c
		JOIN post_category pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("GetCategoriesByPost scan error: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost rows error: %v", err)
	}

	return categories, nil
}

// GetAllCategories returns all available post categories.
// @Summary Get all categories
// @Description Returns a list of all available post categories.
// @Tags Categories
// @Produce json
// @Success 200 {array} string "Categories fetched"
// @Router /api/categories [get]
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	rows, err := sqlite.DB().Query(`SELECT name FROM CATEGORY ORDER BY name ASC`)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch categories", nil)
		return
	}
	defer rows.Close()

	categories := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not scan categories", nil)
			return
		}
		categories = append(categories, name)
	}

	utilities.WriteJSON(w, http.StatusOK, "categories fetched", categories)
}
