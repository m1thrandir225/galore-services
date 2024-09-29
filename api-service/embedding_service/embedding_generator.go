package embedding

import "net/http"

type Generator interface {
	GenerateEmbedding(text string, client *http.Client) ([]float64, error)
}
