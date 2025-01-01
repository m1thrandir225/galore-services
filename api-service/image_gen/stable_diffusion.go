package image_gen

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
)

type StableDiffusionModel string

const (
	StableDiffusionModelUltra StableDiffusionModel = "ultra"
	StableDiffusionModelCore  StableDiffusionModel = "core"
	StableDiffusionModelSD3   StableDiffusionModel = "sd3"
)

type StableDiffusionGenerator struct {
	BaseURL      string
	ApiKey       string
	OutputFormat string
	AspectRatio  string
}

func NewStableDiffusionGenerator(url, apiKey, aspectRatio, outputFormat string) *StableDiffusionGenerator {
	return &StableDiffusionGenerator{
		BaseURL:      url,
		ApiKey:       apiKey,
		OutputFormat: outputFormat,
		AspectRatio:  aspectRatio,
	}
}

func (generator *StableDiffusionGenerator) generateUrlBasedOnModel(model StableDiffusionModel) (string, error) {
	switch model {
	case StableDiffusionModelUltra:
		return fmt.Sprintf("%s/ultra", generator.BaseURL), nil
	case StableDiffusionModelCore:
		return fmt.Sprintf("%s/core", generator.BaseURL), nil
	case StableDiffusionModelSD3:
		return fmt.Sprintf("%s/sd3", generator.BaseURL), nil
	default:
		return "", fmt.Errorf("unknown model: %s", model)
	}
}

func (generator *StableDiffusionGenerator) GenerateImage(prompt string, httpClient *http.Client, model string) (*GeneratedImage, error) {
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

	urlForModel, err := generator.generateUrlBasedOnModel(StableDiffusionModel(model))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlForModel, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generator.ApiKey))
	req.Header.Set("Accept", "image/*")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code with body: %s", string(body))
	}

	imageGenerated = GeneratedImage{
		FileName: uuid.New().String(),
		FileExt:  fmt.Sprintf(".%s", generator.OutputFormat),
		Content:  body,
	}

	return &imageGenerated, nil
}
