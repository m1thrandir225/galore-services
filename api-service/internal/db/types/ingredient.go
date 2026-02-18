// Package types defines JSON types for the database.
package types

type IngredientData struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type IngredientDTO struct {
	Ingredients []IngredientData `json:"ingredients"`
}
