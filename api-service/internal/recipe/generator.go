// Package recipe manages recipe generation
package recipe

import (
	"fmt"
	"strings"

	dto2 "github.com/m1thrandir225/galore-services/internal/dto"
)

type PromptCocktail struct {
	Name         string                   `json:"name"`
	Description  string                   `json:"short_description"`
	Instructions []dto2.PromptInstruction `json:"instructions"`
	Ingredients  []dto2.IngredientData    `json:"ingredients"`
	ImagePrompt  string                   `json:"cocktail_image_prompt"`
}

type PromptRecipe struct {
	Cocktail PromptCocktail `json:"cocktail"`
}

// Generator provides a way to generate a recipe
type Generator interface {
	GenerateRecipe(prompt string) (*PromptCocktail, error)
}

func GeneratePrompt(referenceFlavours, referenceCocktails []string) string {
	return fmt.Sprintf("The user has selected the following flavours as a reference: %s. And as reference cocktails he has selected the following: %s", strings.Join(referenceFlavours, ", "), strings.Join(referenceCocktails, ", "))
}
