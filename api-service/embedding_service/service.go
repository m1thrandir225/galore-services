package embedding

// Description:
// Inteface for a embedding service
type EmbeddingService interface {
	// Description:
	// Generate an embedding from a given text
	GenerateEmbedding(text string) ([]float32, error)
}
