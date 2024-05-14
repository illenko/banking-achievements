package model

import "github.com/google/uuid"

type Achievement struct {
	ID           uuid.UUID
	SettingID    uuid.UUID
	Value        float64
	ValuesHolder []string
}
