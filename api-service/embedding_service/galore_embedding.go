package embedding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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

func NewGaloreEmbeddingService(url, apiKey string) *GaloreEmbeddingService {
	return &GaloreEmbeddingService{
		Url:    url,
		ApiKey: apiKey,
	}
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
	req.Header.Set("X-Api-Key", generator.ApiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

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
