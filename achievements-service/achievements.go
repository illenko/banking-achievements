package main

import (
	"github.com/google/uuid"
)

type achievement struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       int       `json:"value"`
	Goal        int       `json:"goal"`
	Repeatable  bool      `json:"repeatable"`
	Count       int       `json:"count"`
}
