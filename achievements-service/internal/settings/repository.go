package settings

import "github.com/google/uuid"

var settings = []Setting{
	{
		Id:          uuid.MustParse("9579af22-bb86-4e8c-8076-218d1a61b70e"),
		Name:        "Big Spender",
		Description: "Made 3 transactions with amount more than $100",
		Filter: Filter{
			Amount: 100,
		},
		Criteria: Criteria{
			Type:  Count,
			Value: 3,
		},
	},
	{
		Id:          uuid.MustParse("ec42bee0-7d36-48de-bf6f-eb9bfd7d4a6e"),
		Name:        "Coffee Addict",
		Description: "Spent more than $50 on coffee",
		Filter: Filter{
			Categories: &[]string{"coffee"},
			Amount:     50,
		},
		Criteria: Criteria{
			Field: Amount,
			Type:  Sum,
			Value: 50,
		},
		Repeatable: true,
	},
	{
		Id:          uuid.MustParse("0c04db79-397a-44aa-b309-d0a7a368424d"),
		Name:        "Traveller",
		Description: "Made transactions in 5 different countries",
		Filter: Filter{
			Amount: 0,
		},
		Criteria: Criteria{
			Field: Country,
			Type:  Unique,
			Value: 5,
		},
	},
	{
		Id:          uuid.MustParse("a04a3397-075c-496b-a636-04f668116a0f"),
		Name:        "Taxi Lover",
		Description: "Made 5 transactions with taxi category",
		Filter: Filter{
			Categories: &[]string{"taxi"},
		},
		Criteria: Criteria{
			Type:  Count,
			Value: 5,
		},
		Repeatable: true,
	},
}

type Repository interface {
	GetAll() ([]Setting, error)
}

type repository struct {
}

func NewSettingsRepository() Repository {
	return &repository{}
}

func (r *repository) GetAll() ([]Setting, error) {
	return settings, nil
}
