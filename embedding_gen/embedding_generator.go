package embedding

type EmbeddingGenerator interface {
	GenerateEmbedding(text string) ([]float64, error)
}
