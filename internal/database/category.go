package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{
		db: db,
	}
}

func (c *Category) Create(name, description string) (Category, error) {
	id := uuid.NewString()

	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES (?, ?, ?)", id, name, description)
	if err != nil {
		return Category{}, err
	}

	return Category{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}

func (c *Category) GetByID(id string) (Category, error) {
	var category Category

	row := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = ?", id)
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (c *Category) ListAll() ([]Category, error) {
	var categories []Category

	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (c *Category) GetByCourseID(courseID string) (Category, error) {
	var category Category

	row := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = (SELECT category_id FROM courses WHERE id = ?)", courseID)
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}
