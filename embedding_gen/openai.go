package embedding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type OpenAiEmbeddingGenerator struct {
	AuthorizationToken string
	URL                string
	Model              string
}

type OpenAiEmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func (generator *OpenAiEmbeddingGenerator) GenerateEmbedding(text string) ([]float64, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"input": text,
		"model": generator.Model,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", generator.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", generator.AuthorizationToken)

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

	var embeddingResp OpenAiEmbeddingResponse
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		return nil, err
	}
	return embeddingResp.Data[0].Embedding, nil
}
