package main

import (
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	"time"
)

type Field int
type Type int

const (
	Amount Field = iota
	Country
)

const (
	Single Type = iota
	Sum
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
		Description: "Spent more than $100 in a single transaction",
		Filter: achievementFilter{
			Amount: 100,
		},
		Criteria: achievementCriteria{
			Field: Amount,
			Type:  Single,
			Value: 100,
		},
	},
	{
		Id:          uuid.New(),
		Name:        "Coffee Addict",
		Description: "Spent more than $50 on coffee",
		Filter: achievementFilter{
			Categories: &[]string{"Coffee"},
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
		Description: "Made a transactions in 5 different countries",
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
			Categories: &[]string{"Taxi"},
		},
		Criteria: achievementCriteria{
			Type:  Count,
			Value: 5,
		},
	},
}

type achievementValue struct {
	ID        uuid.UUID `json:"id"`
	SettingId uuid.UUID `json:"setting_id"`
	Value     int       `json:"value"`
}

type achievementUniqueValues struct {
	ID        uuid.UUID `json:"id"`
	SettingId uuid.UUID `json:"setting_id"`
	Values    []string  `json:"values"`
}

type achievement struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CriteriaValue int       `json:"criteriaValue"`
	ActualValue   int       `json:"actualValue"`
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

	app.Subscribe("transactions", func(c *gofr.Context) error {

		var data transaction

		err := c.Bind(&data)
		if err != nil {
			c.Logger.Error("Error unmarshalling transaction: ", err)
			return nil
		}

		c.Logger.Info("Consumed transaction", data)

		return nil
	})

	app.Run()
}
