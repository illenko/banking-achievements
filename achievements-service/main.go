package main

import (
	"github.com/illenko/achievements-service/migrations"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	app.Subscribe("transactions", processTransaction)
	app.GET("/achievements", fetchAchievements)

	app.Run()
}
