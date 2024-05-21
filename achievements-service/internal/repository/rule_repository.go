package repository

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/model"
)

var rules = []model.Rule{
	{
		ID:          uuid.MustParse("9579af22-bb86-4e8c-8076-218d1a61b70e"),
		Name:        "🤑 Big Spender",
		Description: "Made 3 transactions with amount more than $100",
		Filter: model.Filter{
			Amount: 100,
		},
		Criteria: model.Criteria{
			Type:  model.Count,
			Value: 3,
		},
	},
	{
		ID:          uuid.MustParse("ec42bee0-7d36-48de-bf6f-eb9bfd7d4a6e"),
		Name:        "☕ Coffee Addict",
		Description: "Spent more than $50 on coffee",
		Filter: model.Filter{
			Categories: &[]string{"coffee"},
			Amount:     50,
		},
		Criteria: model.Criteria{
			Field: model.Amount,
			Type:  model.Sum,
			Value: 50,
		},
		Repeatable: true,
	},
	{
		ID:          uuid.MustParse("0c04db79-397a-44aa-b309-d0a7a368424d"),
		Name:        "🧳 Traveller",
		Description: "Made transactions in 5 different countries",
		Filter: model.Filter{
			Amount: 0,
		},
		Criteria: model.Criteria{
			Field: model.Country,
			Type:  model.Unique,
			Value: 5,
		},
	},
	{
		ID:          uuid.MustParse("a04a3397-075c-496b-a636-04f668116a0f"),
		Name:        "🚕 Taxi Lover",
		Description: "Made 5 transactions with taxi category",
		Filter: model.Filter{
			Categories: &[]string{"taxi"},
		},
		Criteria: model.Criteria{
			Type:  model.Count,
			Value: 5,
		},
		Repeatable: true,
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
