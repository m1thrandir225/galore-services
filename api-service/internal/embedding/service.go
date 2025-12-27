// Package embedding provides a way to return an embedding for a given text
package embedding

type Service interface {
	GenerateEmbedding(text string) ([]float32, error)
}
