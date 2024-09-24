package image_gen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type StableDiffusionGenerator struct {
	Url          string
	ApiKey       string
	OutputFormat string
	AspectRatio  string
	StylePreset  string
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

// GenerateImages TODO: if an error pops up stop the context execution, current implementation will spend credits even if there is an error.
func (generator *StableDiffusionGenerator) GenerateImages(prompts []string, httpClient *http.Client) ([]*GeneratedImage, error) {
	var waitGroup sync.WaitGroup
	imageChannel := make(chan *GeneratedImage, len(prompts))
	errorChannel := make(chan error, len(prompts))

	for _, prompt := range prompts {
		waitGroup.Add(1)

		go func(p string) {
			defer waitGroup.Done()

			image, err := generator.GenerateImage(p, httpClient)
			if err != nil {
				errorChannel <- err
			} else {
				imageChannel <- image
			}
		}(prompt)
	}
	waitGroup.Wait()
	close(imageChannel)
	close(errorChannel)

	var images []*GeneratedImage
	var errorz []error

	for img := range imageChannel {
		images = append(images, img)
	}

	for err := range errorChannel {
		log.Println(err)
		errorz = append(errorz, err)
	}
	if len(errorz) > 0 {
		return nil, errors.New("there was a problem with the generation process, please try again")
	}
	return images, nil
}
