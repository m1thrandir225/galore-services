package embedding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGaloreEmbedding(t *testing.T) {
	embeddingService := GaloreEmbeddingService{
		Url:    "http://localhost:8000/generate-embedding",
		ApiKey: "testing",
	}
	data, err := embeddingService.GenerateEmbedding("Hello World")

	require.NoError(t, err)
	require.NotEmpty(t, data)
	require.Equal(t, len(data), 768)
}
