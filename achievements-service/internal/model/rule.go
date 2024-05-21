package model

import "github.com/google/uuid"

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
	Filter      Filter
	Criteria    Criteria
	Repeatable  bool
}

type Filter struct {
	Categories *[]string
	Amount     float64
}

type Criteria struct {
	Field Field
	Type  Type
	Value int
}
