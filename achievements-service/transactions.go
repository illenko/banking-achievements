package main

import (
	"time"

	"github.com/google/uuid"
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
	achievements, err := fetchAchievementEntities(c)
	if err != nil {
		return err
	}

	return processAchievements(c, t, settings, achievements)
}
