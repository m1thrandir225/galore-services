package image_gen

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
)

type StableDiffusionGenerator struct {
	Url          string
	ApiKey       string
	OutputFormat string
	AspectRatio  string
}

func NewStableDiffusionGenerator(url, apiKey, aspectRatio, outputFormat string) *StableDiffusionGenerator {
	return &StableDiffusionGenerator{
		Url:          url,
		ApiKey:       apiKey,
		OutputFormat: outputFormat,
		AspectRatio:  aspectRatio,
	}
}

func (generator *StableDiffusionGenerator) GenerateImage(prompt string, httpClient *http.Client) (*GeneratedImage, error) {
	buffer := &bytes.Buffer{}
	var imageGenerated GeneratedImage
	mpw := multipart.NewWriter(buffer)

	// 1. Add prompt field to multipart form
	promptField, err := mpw.CreateFormField("prompt")
	if err != nil {
		return nil, err
	}
	_, err = promptField.Write([]byte(prompt))
	if err != nil {
		return nil, err
	}

	// 2. Add aspect ratio to multipart form
	aspectRatioField, err := mpw.CreateFormField("aspect_ratio")
	if err != nil {
		return nil, err
	}
	_, err = aspectRatioField.Write([]byte(generator.AspectRatio))
	if err != nil {
		return nil, err
	}

	// 3. Add output to multipart form
	outputFormatField, err := mpw.CreateFormField("output_format")
	if err != nil {
		return nil, err
	}
	_, err = outputFormatField.Write([]byte(generator.OutputFormat))
	if err != nil {
		return nil, err
	}

	err = mpw.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", generator.Url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generator.ApiKey))
	req.Header.Set("Accept", fmt.Sprintf("image/%s", generator.OutputFormat))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Check HTTP status code before processing body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	imageGenerated = GeneratedImage{
		FileName: uuid.New().String(),
		FileExt:  fmt.Sprintf(".%s", generator.OutputFormat),
		Content:  body,
	}

	return &imageGenerated, nil
}
