package main

import (
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	"time"
)

type transaction struct {
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Merchant string    `json:"merchant"`
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
