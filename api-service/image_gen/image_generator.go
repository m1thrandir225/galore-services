package image_gen

import "net/http"

// Description:
// The generated image result
type GeneratedImage struct {
	FileName string
	FileExt  string
	Content  []byte
}

// Description:
// ImageGenerator interface that specifies the required methods for the service
type ImageGenerator interface {
	GenerateImage(prompt string, httpClient *http.Client, model string) (*GeneratedImage, error)
}
