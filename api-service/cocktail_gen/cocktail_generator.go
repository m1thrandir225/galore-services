package cocktail_gen

import (
	"fmt"
	"github.com/m1thrandir225/galore-services/dto"
	"strings"
)

type PromptCocktail struct {
	Name         string               `json:"name"`
	Description  string               `json:"short_description"`
	Instructions dto.AiInstructionDto `json:"instructions"`
	Ingredients  []dto.IngredientDto  `json:"ingredients"`
	ImagePrompt  string               `json:"cocktail_image_prompt"`
}

type PromptRecipe struct {
	Cocktail PromptCocktail `json:"cocktail"`
}

type CocktailGenerator interface {
	GenerateRecipe(referenceFlavours, referenceCocktails []string) (*PromptRecipe, error)
}

func generatePrompt(referenceFlavours, referenceCocktails []string) string {
	return fmt.Sprintf("The user has selected the following flavours as a reference: %s. And as reference cocktails he has selected the following: %s", strings.Join(referenceFlavours, ", "), strings.Join(referenceCocktails, ", "))
}
