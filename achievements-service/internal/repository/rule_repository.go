package repository

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/internal/operation"
)

var rules = []model.Rule{
	{
		ID:          uuid.MustParse("9579af22-bb86-4e8c-8076-218d1a61b70e"),
		Name:        "🤑 Big Spender",
		Description: "Made 3 transactions with amount more than $100",
		Filters: []model.Filter{
			{
				Field:     "amount",
				Operation: ">",
				Value:     "100",
			},
		},
		Criteria: model.Criteria{
			Operation: operation.Count,
			Value:     3,
		},
	},
	{
		ID:          uuid.MustParse("ec42bee0-7d36-48de-bf6f-eb9bfd7d4a6e"),
		Name:        "☕ Coffee Addict",
		Description: "Spent more than $50 on coffee",
		Filters: []model.Filter{
			{
				Field:     "category",
				Operation: "=",
				Value:     "coffee",
			},
		},
		Criteria: model.Criteria{
			Field:     "amount",
			Operation: operation.Sum,
			Value:     50,
		},
		Repeatable: true,
	},
	{
		ID:          uuid.MustParse("0c04db79-397a-44aa-b309-d0a7a368424d"),
		Name:        "🧳 Traveller",
		Description: "Made transactions in 5 different countries",
		Criteria: model.Criteria{
			Field:     "country",
			Operation: operation.Unique,
			Value:     5,
		},
	},
	{
		ID:          uuid.MustParse("f4b3b3b4-4b3b-4b3b-4b3b-4b3b4b3b4b3b"),
		Name:        "🍔 Foodie",
		Description: "Made 10 transactions with food",
		Filters: []model.Filter{
			{
				Field:     "category",
				Operation: "in",
				Value:     "food,restaurant",
			},
		},
		Criteria: model.Criteria{
			Operation: operation.Count,
			Value:     10,
		},
	},
}

type RuleRepository interface {
	GetAll() ([]model.Rule, error)
}

type ruleRepository struct {
}

func NewRuleRepository() RuleRepository {
	return &ruleRepository{}
}

func (r *ruleRepository) GetAll() ([]model.Rule, error) {
	return rules, nil
}
