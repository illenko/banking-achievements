package main

import (
	"github.com/illenko/achievements-service/internal/migration"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migration.All())

	app.Subscribe("transactions", processTransaction)
	app.GET("/achievements", fetchAchievements)

	app.Run()
}
