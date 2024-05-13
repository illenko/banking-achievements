package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gofr.dev/pkg/gofr"
)

type AchievementRepository interface {
	GetAchievements(c *gofr.Context) (map[uuid.UUID]achievementEntity, error)
	InsertAchievement(c *gofr.Context, a achievementEntity) error
	UpdateAchievement(c *gofr.Context, a achievementEntity) error
}

type achievementRepository struct {
}

func NewAchievementRepository() AchievementRepository {
	return &achievementRepository{}
}

func (r *achievementRepository) GetAchievements(c *gofr.Context) (map[uuid.UUID]achievementEntity, error) {
	rows, err := c.SQL.Query("SELECT id, value, values FROM achievement_entity")
	if err != nil {
		return nil, fmt.Errorf("error querying achievement entities: %w", err)
	}
	defer rows.Close()

	achievements := make(map[uuid.UUID]achievementEntity)
	for rows.Next() {
		var a achievementEntity
		var values []string
		err := rows.Scan(&a.ID, &a.Value, pq.Array(&values))
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		a.Values = values
		achievements[a.ID] = a
	}

	return achievements, nil
}

func (r *achievementRepository) InsertAchievement(c *gofr.Context, a achievementEntity) error {
	_, err := c.SQL.Exec("INSERT INTO achievement_entity (id, value, values) VALUES ($1, $2, $3)", a.ID, a.Value, pq.Array(a.Values))
	if err != nil {
		return fmt.Errorf("error inserting achievement entity: %w", err)
	}

	return nil
}

func (r *achievementRepository) UpdateAchievement(c *gofr.Context, a achievementEntity) error {
	_, err := c.SQL.Exec("UPDATE achievement_entity SET value = $2, values = $3 WHERE id = $1", a.ID, a.Value, pq.Array(a.Values))
	if err != nil {
		return fmt.Errorf("error updating achievement entity: %w", err)
	}

	return nil
}
