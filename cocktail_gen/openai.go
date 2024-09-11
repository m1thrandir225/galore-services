package cocktail_gen

import (
	"bytes"
	"io"
	"net/http"
)

type OpenAIPromptGenerator struct {
	ApiKey      string
	AssistantId string
	ThreadUrl   string
}

func (generator *OpenAIPromptGenerator) GenerateRecipe(referenceFlavours, referenceCocktails []string) (*PromptRecipe, error) {
	/**
	Steps to run
	1. Create a prompt
	2. Create a new thread
	3. Add a message to the thread with the prompt
	4. Wait for the result (i don't know if this is going to be the request or not)
	5. return the prompt recipe.
	*/
	var recipe PromptRecipe
	//1. Create a prompt
	prompt := generatePrompt(referenceFlavours, referenceCocktails)
	httpClient := &http.Client{}
	//2. Create a new thread
	req, err := http.NewRequest("POST", generator.ThreadUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+generator.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}
