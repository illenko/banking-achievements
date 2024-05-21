package processor

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/internal/service"
	"gofr.dev/pkg/gofr"
)

type TransactionProcessor interface {
	Process(c *gofr.Context) error
}

type transactionProcessor struct {
	achievementService service.AchievementService
	rulesService       service.RuleService
}

func NewTransactionProcessor(achievementService service.AchievementService, rulesService service.RuleService) TransactionProcessor {
	return &transactionProcessor{
		achievementService: achievementService,
		rulesService:       rulesService,
	}
}

func (p *transactionProcessor) Process(c *gofr.Context) error {
	var t model.Transaction
	err := c.Bind(&t)
	if err != nil {
		c.Logger.Error("Error unmarshalling transaction: ", err)
		return nil
	}
	c.Logger.Info("Consumed transaction", t)

	achievementrules, err := p.rulesService.GetAll()

	if err != nil {
		return err
	}

	achievements, err := p.achievementService.GetAll(c)
	if err != nil {
		return err
	}

	return p.processAchievements(c, t, filterAchievements(t, achievementrules), achievements)
}

func (p *transactionProcessor) processAchievements(c *gofr.Context, t model.Transaction, rules []model.Rule, achievements []model.Achievement) error {
	achievementMap := make(map[uuid.UUID]model.Achievement)
	for _, a := range achievements {
		achievementMap[a.RuleID] = a
	}

	for _, rule := range rules {
		a, ok := achievementMap[rule.ID]
		if !ok {
			a = newAchievement(rule.ID)
		}

		err := p.processAchievement(c, t, rule, a, ok)
		if err != nil {
			return fmt.Errorf("error processing achievement: %w", err)
		}
	}

	return nil
}

func newAchievement(ruleId uuid.UUID) model.Achievement {
	return model.Achievement{
		ID:           uuid.New(),
		RuleID:       ruleId,
		Value:        0,
		ValuesHolder: []string{},
	}
}

func (p *transactionProcessor) processAchievement(c *gofr.Context, t model.Transaction, s model.Rule, a model.Achievement, existing bool) error {
	switch s.Criteria.Type {
	case model.Count:
		a.Value++
	case model.Sum:
		a.Value += t.Amount
	case model.Unique:
		value, err := getValue(s.Criteria, t)
		if err != nil {
			return fmt.Errorf("error getting value: %w", err)
		}
		if !slices.Contains(a.ValuesHolder, value) {
			a.ValuesHolder = append(a.ValuesHolder, value)
			a.Value++
		}
	}

	if existing {
		err := p.achievementService.UpdateAchievement(c, a)
		if err != nil {
			return fmt.Errorf("error updating achievement entity: %w", err)
		}
	} else {
		err := p.achievementService.InsertAchievement(c, a)
		if err != nil {
			return fmt.Errorf("error creating achievement entity: %w", err)
		}
	}

	return nil
}

func getValue(criteria model.Criteria, t model.Transaction) (string, error) {
	if criteria.Field == model.Country {
		return t.Country, nil
	} else if criteria.Field == model.Category {
		return t.Category, nil
	} else {
		return "", fmt.Errorf("unknown field: %v", criteria.Field)
	}
}

func filterAchievements(t model.Transaction, s []model.Rule) []model.Rule {
	var filteredAchievements []model.Rule

	for _, rule := range s {
		if meetsCriteria(t, rule.Filter) {
			filteredAchievements = append(filteredAchievements, rule)
		}
	}

	return filteredAchievements
}

func meetsCriteria(t model.Transaction, filter model.Filter) bool {
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
	} else {
		categoryMatched = true
	}

	return categoryMatched
}
