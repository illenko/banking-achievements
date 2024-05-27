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
	var transaction model.Transaction
	err := c.Bind(&transaction)
	if err != nil {
		c.Logger.Error("Error unmarshalling transaction: ", err)
		return nil
	}
	c.Logger.Info("Consumed transaction", transaction)

	rules, err := p.rulesService.GetAll()

	if err != nil {
		return err
	}

	achievements, err := p.achievementService.GetAll(c)
	if err != nil {
		return err
	}

	return p.processAchievements(c, transaction, filterAchievements(transaction, rules), achievements)
}

func (p *transactionProcessor) processAchievements(c *gofr.Context, transaction model.Transaction, rules []model.Rule, achievements []model.Achievement) error {
	achievementMap := make(map[uuid.UUID]model.Achievement)
	for _, a := range achievements {
		achievementMap[a.RuleID] = a
	}

	for _, rule := range rules {
		achievement, ok := achievementMap[rule.ID]
		if !ok {
			achievement = newAchievement(rule.ID)
		}

		err := p.processAchievement(c, transaction, rule, achievement, ok)
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

func (p *transactionProcessor) processAchievement(c *gofr.Context, transaction model.Transaction, rule model.Rule, achievement model.Achievement, existing bool) error {
	switch rule.Criteria.Type {
	case model.Count:
		achievement.Value++
	case model.Sum:
		achievement.Value += transaction.Amount
	case model.Unique:
		value, err := getValue(rule.Criteria, transaction)
		if err != nil {
			return fmt.Errorf("error getting value: %w", err)
		}
		if !slices.Contains(achievement.ValuesHolder, value) {
			achievement.ValuesHolder = append(achievement.ValuesHolder, value)
			achievement.Value++
		}
	}

	if existing {
		err := p.achievementService.UpdateAchievement(c, achievement)
		if err != nil {
			return fmt.Errorf("error updating achievement entity: %w", err)
		}
	} else {
		err := p.achievementService.InsertAchievement(c, achievement)
		if err != nil {
			return fmt.Errorf("error creating achievement entity: %w", err)
		}
	}

	return nil
}

func getValue(criteria model.Criteria, transaction model.Transaction) (string, error) {
	if criteria.Field == model.Country {
		return transaction.Country, nil
	} else if criteria.Field == model.Category {
		return transaction.Category, nil
	} else {
		return "", fmt.Errorf("unknown field: %v", criteria.Field)
	}
}

func filterAchievements(t model.Transaction, rules []model.Rule) []model.Rule {
	var filteredAchievements []model.Rule

	for _, rule := range rules {
		if meetsCriteria(t, rule.Filter) {
			filteredAchievements = append(filteredAchievements, rule)
		}
	}

	return filteredAchievements
}

func meetsCriteria(transaction model.Transaction, filter model.Filter) bool {
	if filter.Amount > 0 && transaction.Amount < filter.Amount {
		return false
	}

	categoryMatched := false
	if filter.Categories != nil {
		for _, category := range *filter.Categories {
			if category == transaction.Category {
				categoryMatched = true
				break
			}
		}
	} else {
		categoryMatched = true
	}

	return categoryMatched
}
