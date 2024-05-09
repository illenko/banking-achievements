package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gofr.dev/pkg/gofr"
)

type transaction struct {
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Country  string    `json:"country"`
}

func processTransaction(c *gofr.Context) error {
	var t transaction
	err := c.Bind(&t)
	if err != nil {
		c.Logger.Error("Error unmarshalling transaction: ", err)
		return nil
	}
	c.Logger.Info("Consumed transaction", t)

	settings := filterAchievements(t, achievementSettings)
	rows, err := c.SQL.Query("SELECT id, value, values FROM achievement_entity")
	if err != nil {
		return fmt.Errorf("error querying achievement entities: %w", err)
	}
	defer rows.Close()

	achievements := make(map[uuid.UUID]achievementEntity)
	for rows.Next() {
		var a achievementEntity
		var values []string
		err := rows.Scan(&a.ID, &a.Value, pq.Array(&values))
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		a.Values = values
		achievements[a.ID] = a
	}

	for _, setting := range settings {
		err := processAchievement(c, t, setting, achievements)
		if err != nil {
			return fmt.Errorf("error processing achievement: %w", err)
		}
	}

	return nil
}
