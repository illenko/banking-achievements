package achievement

import "github.com/google/uuid"

type Achievement struct {
	ID        uuid.UUID
	SettingID uuid.UUID
	Value     float64
	Values    []string
}
