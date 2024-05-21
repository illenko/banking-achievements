package repository

import (
	"github.com/illenko/achievements-service/internal/model"
	"github.com/lib/pq"
	"gofr.dev/pkg/gofr"
)

type AchievementRepository interface {
	GetAll(c *gofr.Context) ([]model.Achievement, error)
	Insert(c *gofr.Context, a model.Achievement) error
	Update(c *gofr.Context, a model.Achievement) error
}

type achievementRepository struct {
}

func NewAchievementRepository() AchievementRepository {
	return &achievementRepository{}
}

func (r *achievementRepository) GetAll(c *gofr.Context) ([]model.Achievement, error) {
	rows, err := c.SQL.Query("SELECT id, rule_id, value, values_holder FROM achievement")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err := rows.Scan(&a.ID, &a.RuleID, &a.Value, pq.Array(&a.ValuesHolder))
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}

	return achievements, nil
}

func (r *achievementRepository) Insert(c *gofr.Context, a model.Achievement) error {
	_, err := c.SQL.Exec("INSERT INTO achievement (id, rule_id, value, values_holder) VALUES ($1, $2, $3, $4)", a.ID, a.RuleID, a.Value, pq.Array(a.ValuesHolder))
	if err != nil {
		return err
	}

	return nil
}

func (r *achievementRepository) Update(c *gofr.Context, a model.Achievement) error {
	_, err := c.SQL.Exec("UPDATE achievement SET rule_id = $1, value = $2, values_holder = $3 WHERE id = $4", a.RuleID, a.Value, pq.Array(a.ValuesHolder), a.ID)
	if err != nil {
		return err
	}

	return nil
}
