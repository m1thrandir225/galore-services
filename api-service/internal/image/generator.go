// Package image provides a way to generate images from a given prompt
package image

import "net/http"

type GeneratedImage struct {
	FileName string
	FileExt  string
	Content  []byte
}

// Generator handles image generation from a prompt
type Generator interface {
	GenerateImage(prompt string, httpClient *http.Client, model string) (*GeneratedImage, error)
}
