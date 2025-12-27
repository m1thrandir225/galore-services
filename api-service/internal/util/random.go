package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/m1thrandir225/galore-services/internal/dto"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt returns a random integer based on a min and max value.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFloatArray returns an array of length n containing random float integers given a max and minx value.
func RandomFloatArray(min, max float32, n int) []float32 {
	res := make([]float32, n)
	for i := range res {
		res[i] = min + rand.Float32()*(max-min)
	}
	return res
}

// RandomString returns a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomBool returns a random true or false value.
func RandomBool() bool {
	return rand.Uint64()%2 == 1
}

// RandomEmail returns a randomly generated email.
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

// RandomDate returns a random formatted string date based from 1965 up until today's date
// FIXME: implement now date
func RandomDate() string {
	minDate := time.Date(1965, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Unix()

	delta := maxDate - minDate

	sec := rand.Int63n(delta) + minDate

	dateWithTime := time.Unix(sec, 0)

	return fmt.Sprintf("%d-%02d-%02d", dateWithTime.Year(), dateWithTime.Month(), dateWithTime.Day())
}

// RandomIngredients returns an array of length n containing dto.Ingredients
func RandomIngredients(n int) []dto.IngredientData {
	ingredients := make([]dto.IngredientData, n)

	for i := range ingredients {
		newIngredient := dto.IngredientData{
			Name:   RandomString(10),
			Amount: RandomString(4),
		}
		ingredients[i] = newIngredient
	}
	return ingredients
}

// RandomStringArray returns a random array of length n given with a random string of length m.
func RandomStringArray(n, m int) []string {
	strArr := make([]string, n)
	for i := range strArr {
		strArr[i] = RandomString(m)
	}
	return strArr
}
