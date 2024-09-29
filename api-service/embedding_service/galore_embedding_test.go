package embedding

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGaloreEmbedding(t *testing.T) {
	httpClient := &http.Client{}

	embeddingService := GaloreEmbeddingService{
		Url: "http://localhost:8000/generate-embedding",
	}

	data, err := embeddingService.GenerateEmbedding("Hello World", httpClient)

	require.NoError(t, err)
	require.NotEmpty(t, data)
}
