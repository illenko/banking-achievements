package model

import "github.com/google/uuid"

type Achievement struct {
	ID           uuid.UUID
	RuleID       uuid.UUID
	Value        float64
	ValuesHolder []string
}
