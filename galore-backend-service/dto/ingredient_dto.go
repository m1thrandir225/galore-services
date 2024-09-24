package dto

type IngredientData struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type IngredientDto struct {
	Ingredients []IngredientData `json:"ingredients"`
}
