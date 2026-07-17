package repository

import (
	"database/sql"
	"fmt"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAll() ([]Category, error) {
	rows, err := r.DB.Query(`
		SELECT id, name
		FROM CATEGORY
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) GetByName(name string) (*Category, error) {
	var c Category

	err := r.DB.QueryRow(`
		SELECT id, name
		FROM CATEGORY
		WHERE name = ?
	`, name).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetIDsByNames resolves a list of category names to their IDs. It returns an
// error naming the first category that doesn't exist, so callers can validate
// user input before inserting.
func (r *CategoryRepository) GetIDsByNames(names []string) ([]int, error) {
	ids := make([]int, 0, len(names))

	for _, name := range names {
		c, err := r.GetByName(name)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid category: %s", name)
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, c.ID)
	}

	return ids, nil
}

func (r *CategoryRepository) GetByPost(postID int) ([]Category, error) {
	rows, err := r.DB.Query(`
		SELECT c.id, c.name
		FROM CATEGORY c
		JOIN POST_CATEGORY pc
		ON pc.category_id = c.id
		WHERE pc.post_id = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

// GetNamesByPost returns just the category names for a post, ordered by name.
// This replaces the old package-level GetCategoriesByPost helper that reached
// into a global db.Database connection directly.
func (r *CategoryRepository) GetNamesByPost(postID int) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT c.name
		FROM CATEGORY c
		JOIN POST_CATEGORY pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("GetNamesByPost error: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("GetNamesByPost scan error: %v", err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetNamesByPost rows error: %v", err)
	}

	return names, nil
}
