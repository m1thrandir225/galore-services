package util

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type instructionWithImage struct {
	Instruction string `json:"instruction"`
	Image       string `json:"instruction_image"`
}

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
func RandomIngredients() []byte {
	var ingredients []ingredient

	for i := 0; i < 10; i++ {
		new_ingredient := ingredient{
			Name:   RandomString(10),
			Amount: RandomString(4),
		}

		ingredients = append(ingredients, new_ingredient)
	}

	b, err := json.Marshal(ingredients)

	if err != nil {
		log.Fatal("There was a problem encoding json: ", err)
	}

	return b

}

/*
* Generate a random array of instructions and return a json
 */
func RandomInstructions() []byte {
	var instructions []instructionWithImage

	for i := 0; i < 10; i++ {
		new_instruction := instructionWithImage{
			Image:       RandomString(12),
			Instruction: RandomString(64),
		}

		instructions = append(instructions, new_instruction)
	}

	b, err := json.Marshal(instructions)

	if err != nil {
		log.Fatal("There was a problem encoding json: ", err)
	}

	return b
}
