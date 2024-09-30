package embedding

import "net/http"

type EmbeddingService interface {
	GenerateEmbedding(text string, client *http.Client) ([]float32, error)
}
