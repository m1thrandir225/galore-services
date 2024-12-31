package image_gen

import "net/http"

type GeneratedImage struct {
	FileName string
	FileExt  string
	Content  []byte
}

type ImageGenerator interface {
	GenerateImage(prompt string, httpClient *http.Client) (*GeneratedImage, error)
}
