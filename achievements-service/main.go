package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/migrations"
	"github.com/lib/pq"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
	"slices"
	"time"
)

type Field int
type Type int

const (
	Amount Field = iota
	Country
	Category
)

const (
	Sum Type = iota
	Unique
	Count
)

type achievementSetting struct {
	Id          uuid.UUID
	Name        string
	Description string
	Filter      achievementFilter
	Criteria    achievementCriteria
}

type achievementFilter struct {
	Categories *[]string
	Amount     float64
}

type achievementCriteria struct {
	Field Field
	Type  Type
	Value int
}

var achievementSettings = []achievementSetting{
	{
		Id:          uuid.New(),
		Name:        "Big Spender",
		Description: "Made 3 transactions with amount more than $100",
		Filter: achievementFilter{
			Amount: 100,
		},
		Criteria: achievementCriteria{
			Type:  Count,
			Value: 3,
		},
	},
	{
		Id:          uuid.New(),
		Name:        "Coffee Addict",
		Description: "Spent more than $50 on coffee",
		Filter: achievementFilter{
			Categories: &[]string{"coffee"},
			Amount:     50,
		},
		Criteria: achievementCriteria{
			Field: Amount,
			Type:  Sum,
			Value: 50,
		},
	},
	{
		Id:          uuid.New(),
		Name:        "Traveller",
		Description: "Made transactions in 5 different countries",
		Filter: achievementFilter{
			Amount: 0,
		},
		Criteria: achievementCriteria{
			Field: Country,
			Type:  Unique,
			Value: 5,
		},
	},
	{
		Id:          uuid.New(),
		Name:        "Taxi Lover",
		Description: "Made 5 transactions with taxi category",
		Filter: achievementFilter{
			Categories: &[]string{"taxi"},
		},
		Criteria: achievementCriteria{
			Type:  Count,
			Value: 5,
		},
	},
}

type achievementEntity struct {
	ID     uuid.UUID
	Value  float64
	Values []string
}

type achievement struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       int       `json:"value"`
	Goal        int       `json:"goal"`
}

type transaction struct {
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Country  string    `json:"country"`
}

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	app.Subscribe("transactions", func(c *gofr.Context) error {
		var data transaction
		err := c.Bind(&data)
		if err != nil {
			c.Logger.Error("Error unmarshalling transaction: ", err)
			return nil
		}
		c.Logger.Info("Consumed transaction", data)
		return processTransaction(c, data)
	})

	app.GET("/achievements", fetchAchievements)

	app.Run()
}

func fetchAchievements(ctx *gofr.Context) (interface{}, error) {
	var achievements []achievement
	rows, err := ctx.SQL.Query("SELECT id, value, values FROM achievement_entity")
	if err != nil {
		return nil, fmt.Errorf("error querying achievement entities: %w", err)
	}
	defer rows.Close()

	achievementEntities := make(map[uuid.UUID]achievementEntity)
	for rows.Next() {
		var a achievementEntity
		var values []string
		err := rows.Scan(&a.ID, &a.Value, pq.Array(&values))
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		a.Values = values
		achievementEntities[a.ID] = a
	}

	for _, setting := range achievementSettings {
		a, ok := achievementEntities[setting.Id]
		if !ok {
			a = achievementEntity{
				ID:     setting.Id,
				Value:  0,
				Values: []string{},
			}
		}

		achievements = append(achievements, achievement{
			ID:          setting.Id,
			Name:        setting.Name,
			Description: setting.Description,
			Value:       int(a.Value),
			Goal:        setting.Criteria.Value,
		})
	}

	return response.Raw{Data: achievements}, nil
}

func processTransaction(c *gofr.Context, t transaction) error {
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

func processAchievement(c *gofr.Context, t transaction, s achievementSetting, achievements map[uuid.UUID]achievementEntity) error {
	a, ok := achievements[s.Id]
	if !ok {
		a = newAchievementEntity(s.Id)
	}

	switch s.Criteria.Type {
	case Count:
		a.Value++
	case Sum:
		a.Value += t.Amount
	case Unique:
		value, err := getValue(s.Criteria, t)
		if err != nil {
			return fmt.Errorf("error getting value: %w", err)
		}
		if !slices.Contains(a.Values, value) {
			a.Values = append(a.Values, value)
			a.Value++
		}
	}

	if ok {
		err := updateAchievementEntity(c, a)
		if err != nil {
			return fmt.Errorf("error updating achievement entity: %w", err)
		}
	} else {
		err := createAchievementEntity(c, a)
		if err != nil {
			return fmt.Errorf("error creating achievement entity: %w", err)
		}
	}

	return nil
}

func newAchievementEntity(id uuid.UUID) achievementEntity {
	return achievementEntity{
		ID:     id,
		Value:  0,
		Values: []string{},
	}
}

func createAchievementEntity(c *gofr.Context, a achievementEntity) error {
	_, err := c.SQL.Exec("INSERT INTO achievement_entity (id, value, values) VALUES ($1, $2, $3)", a.ID, a.Value, pq.Array(a.Values))
	if err != nil {
		return fmt.Errorf("error executing insert: %w", err)
	}
	return nil
}

func updateAchievementEntity(c *gofr.Context, a achievementEntity) error {
	_, err := c.SQL.Exec("UPDATE achievement_entity SET value = $1, values = $2 WHERE id = $3", a.Value, pq.Array(a.Values), a.ID)
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}
	return nil
}

func getValue(criteria achievementCriteria, t transaction) (string, error) {
	if criteria.Field == Country {
		return t.Country, nil
	} else if criteria.Field == Category {
		return t.Category, nil
	} else {
		return "", fmt.Errorf("unknown field: %v", criteria.Field)
	}
}

func filterAchievements(t transaction, settings []achievementSetting) []achievementSetting {
	var filteredAchievements []achievementSetting

	for _, setting := range settings {
		if meetsCriteria(t, setting.Filter) {
			filteredAchievements = append(filteredAchievements, setting)
		}
	}

	return filteredAchievements
}

func meetsCriteria(t transaction, filter achievementFilter) bool {
	if filter.Amount > 0 && t.Amount < filter.Amount {
		return false
	}

	categoryMatched := false
	if filter.Categories != nil {
		for _, category := range *filter.Categories {
			if category == t.Category {
				categoryMatched = true
				break
			}
		}
	}

	return categoryMatched
}
