package categorizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type GaloreCategorizer struct {
	Url    string
	ApiKey string
}

type GaloreCategorizeRequest struct {
	Cocktail db.Cocktail `json:"cocktail"`
}

func NewGaloreCategorizer(url, apiKey string) *GaloreCategorizer {
	return &GaloreCategorizer{
		Url:    url,
		ApiKey: apiKey,
	}
}

/**
* Calls the categorizer to categorize the current cocktail
 */
func (categorizer *GaloreCategorizer) CategorizeCocktail(cocktail db.Cocktail) error {
	requestJson, err := json.Marshal(cocktail)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		categorizer.Url,
		bytes.NewBuffer(requestJson),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", categorizer.ApiKey)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	defer response.Body.Close()

	return nil
}
