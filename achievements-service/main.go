package main

import (
	"github.com/illenko/achievements-service/internal/handler"
	"github.com/illenko/achievements-service/internal/mapper"
	"github.com/illenko/achievements-service/internal/migration"
	"github.com/illenko/achievements-service/internal/processor"
	"github.com/illenko/achievements-service/internal/repository"
	"github.com/illenko/achievements-service/internal/service"

	"gofr.dev/pkg/gofr"
)

func main() {

	ruleRepository := repository.NewRuleRepository()
	achievementRepository := repository.NewAchievementRepository()
	achievementMapper := mapper.NewAchievementMapper()

	achievementService := service.NewAchievementService(achievementRepository, ruleRepository, achievementMapper)
	ruleService := service.NewRuleService(ruleRepository)

	achievementHandler := handler.NewAchievementHandler(achievementService)
	transactionProcessor := processor.NewTransactionProcessor(achievementService, ruleService)

	app := gofr.New()

	app.Migrate(migration.All())

	app.Subscribe("transactions", transactionProcessor.Process)
	app.GET("/achievements", achievementHandler.GetAll)

	app.Run()
}
