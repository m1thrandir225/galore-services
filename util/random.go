package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	dto "github.com/m1thrandir225/galore-services/dto"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
func RandomBool() bool {
	return rand.Uint64()%2 == 1
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

func RandomDate() string {
	min := time.Date(1965, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()

	delta := max - min

	sec := rand.Int63n(delta) + min

	dateWithTime := time.Unix(sec, 0)

	return fmt.Sprintf("%d-%02d-%02d", dateWithTime.Year(), dateWithTime.Month(), dateWithTime.Day())
}

/**
* Generate a random array of ingredients and return a json
 */
func RandomIngredients() []dto.IngredientData {
	var ingredients []dto.IngredientData

	for i := 0; i < 10; i++ {
		new_ingredient := dto.IngredientData{
			Name:   RandomString(10),
			Amount: RandomString(4),
		}

		ingredients = append(ingredients, new_ingredient)
	}

	return ingredients

}

/*
* Generate a random array of instructions and return a json
 */
func RandomAiInstructions() []dto.AiInstructionData {
	var instructions []dto.AiInstructionData

	for i := 0; i < 10; i++ {
		new_instruction := dto.AiInstructionData{
			InstructionImage: RandomString(12),
			Instruction:      RandomString(64),
		}

		instructions = append(instructions, new_instruction)
	}

	return instructions
}

func RandomInstructions() []dto.InstructionData {
	var instructions []dto.InstructionData

	for i := 0; i < 10; i++ {
		new_instruction := dto.InstructionData{
			Instruction: RandomString(64),
		}

		instructions = append(instructions, new_instruction)
	}

	return instructions
}
