package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type EmbeddingService interface {
	GenerateEmbedding(text string) ([]float32, error)
}

type GaloreEmbeddingService struct {
	Url    string
	ApiKey string
}

type GaloreEmbeddingServiceRequest struct {
	Text string `json:"text" binding:"required"`
}

type GaloreEmbeddingServiceResponse struct {
	Embedding [][]float32 `json:"embedding" binding:"required"`
}

// NewGaloreEmbeddingService returns an instance of the GaloreEmbeddingService
func NewGaloreEmbeddingService(url, apiKey string) (*GaloreEmbeddingService, error) {
	if len(url) == 0 {
		return nil, errors.New("url can't be empty")
	}

	if len(apiKey) == 0 {
		return nil, errors.New("api key can't be empty")
	}

	return &GaloreEmbeddingService{
		Url:    url,
		ApiKey: apiKey,
	}, nil
}

func (generator *GaloreEmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	request := GaloreEmbeddingServiceRequest{
		Text: text,
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", generator.Url, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", generator.ApiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embeddingResponse GaloreEmbeddingServiceResponse

	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	return embeddingResponse.Embedding[0], nil
}
