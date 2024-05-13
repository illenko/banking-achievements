package achievement

import (
	"gofr.dev/pkg/gofr"
)

type Repository interface {
	GetAll(c *gofr.Context) ([]Achievement, error)
	Insert(c *gofr.Context, a Achievement) error
	Update(c *gofr.Context, a Achievement) error
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetAll(c *gofr.Context) ([]Achievement, error) {
	rows, err := c.SQL.Query("SELECT id, value, values FROM achievement_entity")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []Achievement
	for rows.Next() {
		var a Achievement
		err := rows.Scan(&a.ID, &a.Value, &a.Values)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}

	return achievements, nil
}

func (r *repository) Insert(c *gofr.Context, a Achievement) error {
	_, err := c.SQL.Exec("INSERT INTO achievement_entity (id, value, values) VALUES ($1, $2, $3)", a.ID, a.Value, a.Values)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(c *gofr.Context, a Achievement) error {
	_, err := c.SQL.Exec("UPDATE achievement_entity SET value = $2, values = $3 WHERE id = $1", a.ID, a.Value, a.Values)
	if err != nil {
		return err
	}

	return nil
}
