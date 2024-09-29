package embedding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type GaloreEmbeddingService struct {
	Url string
}

type GaloreEmbeddingServiceRequest struct {
	Text string `json:"text" binding:"required"`
}

type GaloreEmbeddingServiceResponse struct {
	Embedding [][]float64 `json:"embedding" binding:"required"`
}

func (generator *GaloreEmbeddingService) GenerateEmbedding(text string, client *http.Client) ([]float64, error) {
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
