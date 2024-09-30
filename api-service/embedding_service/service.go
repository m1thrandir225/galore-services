package embedding

type EmbeddingService interface {
	GenerateEmbedding(text string) ([]float32, error)
}
