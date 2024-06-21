package processor

import (
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
	var message float64
	err := c.Bind(&message)
	if err != nil {
		c.Logger.Error("Error unmarshalling message: ", err)
		return nil
	}

	c.Logger.Info("Consumed message: ", message)

	return nil
}
