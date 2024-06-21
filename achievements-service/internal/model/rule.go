package model

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/operation"
)

type Field int
type Type int

const (
	Amount Field = iota
	Country
	Category
)

const (
	Sum Type = iota
	Unique
	Count
)

type Rule struct {
	ID          uuid.UUID
	Name        string
	Description string
	Filters     []Filter
	Criteria    Criteria
	Repeatable  bool
}

type Filter struct {
	Field     string
	Operation string
	Value     string
}

type Criteria struct {
	Field     string
	Operation operation.Operation
	Value     int
}
