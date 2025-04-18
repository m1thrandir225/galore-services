package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	dto "github.com/m1thrandir225/galore-services/dto"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Description:
// Return a random integer based on a min and max value.
//
// Parameters:
// min: int64,
// max: int64
//
// Return:
// int64
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Description:
// Rerturn a random array of float integers based where each is based on a min,
// max value.
//
// Parameters:
// min: float32,
// max: float32,
// n: int
//
// Return:
// []float32
func RandomFloatArray(min, max float32, n int) []float32 {
	res := make([]float32, n)
	for i := range res {
		res[i] = min + rand.Float32()*(max-min)
	}
	return res
}

// Description:
// Return a random string based on a given length.
//
// Parameters:
// n: int
//
// Return:
// string
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Description:
// Return a random boolean
//
// Return:
// bool
func RandomBool() bool {
	return rand.Uint64()%2 == 1
}

// Description:
// Return a randomly generated email
//
// Return:
// string
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

// Description:
// Return a random date based from 1965 up until 2024//FIXME: implement now date
//
// Return:
// string
func RandomDate() string {
	min := time.Date(1965, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()

	delta := max - min

	sec := rand.Int63n(delta) + min

	dateWithTime := time.Unix(sec, 0)

	return fmt.Sprintf("%d-%02d-%02d", dateWithTime.Year(), dateWithTime.Month(), dateWithTime.Day())
}

// Description:
// Return a random array of the data-transfer object Ingreident Data
//
// Return:
// []dto.IngredientData
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

// Description:
// Return a random array of the data-transfer object AiInstructionData
//
// Return:
// []dto.AiInstructionData
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

// Description:
// Return a random array of instruction
//
// Return:
// []string
func RandomInstructions() []string {
	var instructions []string

	for i := 0; i < 10; i++ {
		string := RandomString(32)
		instructions = append(instructions, string)
	}

	return instructions
}
