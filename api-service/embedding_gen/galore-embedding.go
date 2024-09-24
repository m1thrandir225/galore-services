package embedding

import "net/http"

type GaloreEmbeddingService struct {
	Url string
}

type GaloreEmbeddingServiceResponse struct {
	Embedding []float64 `json:"embedding" binding:"required"`
}

// TODO: implement
func (generator *GaloreEmbeddingService) GenerateEmbedding(text string, client *http.Client) ([]float64, error) {
	return nil, nil
}
