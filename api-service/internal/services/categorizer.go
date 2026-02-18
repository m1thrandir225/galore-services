package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type CategorizerService interface {
	Categorize(cocktail db.Cocktail) error
}

// GaloreCategorizer is an internal implementation of the Service with the categorizer microservice.
type GaloreCategorizer struct {
	Url    string
	ApiKey string
}

type GaloreCategorizeRequest struct {
	Cocktail db.Cocktail `json:"cocktail"`
}

// NewGaloreCategorizer returns a GaloreCategorizer instance
func NewGaloreCategorizer(url, apiKey string) (*GaloreCategorizer, error) {
	if len(url) == 0 {
		return nil, errors.New("url can't be empty")
	}
	if len(apiKey) == 0 {
		return nil, errors.New("api key can't be empty")
	}
	return &GaloreCategorizer{
		Url:    url,
		ApiKey: apiKey,
	}, nil
}

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
	if err != nil {
		return err
	}

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

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	return nil
}
