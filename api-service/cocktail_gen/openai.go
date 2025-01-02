package cocktail_gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

type OpenAIPromptGenerator struct {
	ApiKey      string
	AssistantId string
	ThreadUrl   string
}

func NewOpenAIPromptGenerator(apiKey, assistantId, threadUrl string) *OpenAIPromptGenerator {
	return &OpenAIPromptGenerator{
		ApiKey:      apiKey,
		AssistantId: assistantId,
		ThreadUrl:   threadUrl,
	}
}

type ThreadMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Thread struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int    `json:"created_at"`
	Metadata  struct {
	} `json:"metadata"`
	ToolResources struct {
	} `json:"tool_resources"`
}

type RunPayload struct {
	AssistantId string `json:"assistant_id"`
}

type Run struct {
	Id                string      `json:"id"`
	Object            string      `json:"object"`
	CreatedAt         int         `json:"created_at"`
	AssistantId       string      `json:"assistant_id"`
	ThreadId          string      `json:"thread_id"`
	Status            string      `json:"status"`
	StartedAt         int         `json:"started_at"`
	ExpiresAt         interface{} `json:"expires_at"`
	CancelledAt       interface{} `json:"cancelled_at"`
	FailedAt          interface{} `json:"failed_at"`
	CompletedAt       int         `json:"completed_at"`
	LastError         interface{} `json:"last_error"`
	Model             string      `json:"model"`
	Instructions      interface{} `json:"instructions"`
	IncompleteDetails interface{} `json:"incomplete_details"`
	Tools             []struct {
		Type string `json:"type"`
	} `json:"tools"`
	Metadata struct {
	} `json:"metadata"`
	Usage               interface{} `json:"usage"`
	Temperature         float64     `json:"temperature"`
	TopP                float64     `json:"top_p"`
	MaxPromptTokens     int         `json:"max_prompt_tokens"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	TruncationStrategy  struct {
		Type         string      `json:"type"`
		LastMessages interface{} `json:"last_messages"`
	} `json:"truncation_strategy"`
	ResponseFormat    interface{} `json:"response_format"`
	ToolChoice        interface{} `json:"tool_choice"`
	ParallelToolCalls bool        `json:"parallel_tool_calls"`
}

type Message struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int    `json:"created_at"`
	ThreadId  string `json:"thread_id"`
	Role      string `json:"role"`
	Content   []struct {
		Type string `json:"type"`
		Text struct {
			Value       string        `json:"value"`
			Annotations []interface{} `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
	AssistantId string        `json:"assistant_id"`
	RunId       string        `json:"run_id"`
	Attachments []interface{} `json:"attachments"`
	Metadata    struct {
	} `json:"metadata"`
}

type ListMessagesResponse struct {
	Object string    `json:"object"`
	Data   []Message `json:"data"`
}

func (generator *OpenAIPromptGenerator) createThread(httpClient http.Client) (*Thread, error) {
	req, err := http.NewRequest("POST", generator.ThreadUrl, nil)
	if err != nil {
		return nil, errors.New("there was a problem creating the request: " + err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+generator.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	threadResponse, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New("there was a problem making the request: " + err.Error())
	}
	defer threadResponse.Body.Close()

	threadBody, err := io.ReadAll(threadResponse.Body)
	if err != nil {
		return nil, errors.New("there was a problem reading body: " + err.Error())
	}

	var threadMessage Thread

	err = json.Unmarshal(threadBody, &threadMessage)
	if err != nil {
		return nil, errors.New("there was a problem parsing the json: " + err.Error())
	}

	return &threadMessage, nil
}

func (thread *Thread) addMessageToThread(httpClient http.Client, message, threadUrl, apiKey string) error {
	payload := ThreadMessage{
		Role:    "user",
		Content: message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.New("there was a problem marshalling the json: " + err.Error())
	}

	messageThreadUrl := fmt.Sprintf("%s/%s/messages", threadUrl, thread.Id)

	messageReq, err := http.NewRequest("POST", messageThreadUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.New("there was a problem creating the request: " + err.Error())
	}

	messageReq.Header.Add("Authorization", "Bearer "+apiKey)
	messageReq.Header.Add("Content-Type", "application/json")
	messageReq.Header.Add("OpenAI-Beta", "assistants=v2")

	messageResponse, err := httpClient.Do(messageReq)
	if err != nil {
		return errors.New("there was a problem making the request: " + err.Error())
	}

	defer messageResponse.Body.Close()

	if messageResponse.StatusCode != 200 {
		return errors.New("there was a problem making the request: " + messageResponse.Status)
	}
	return nil
}

func (generator *OpenAIPromptGenerator) createRun(httpClient http.Client, thread *Thread) (*Run, error) {
	payload := RunPayload{
		AssistantId: generator.AssistantId,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("there was a problem marshalling the json: " + err.Error())
	}
	runUrl := fmt.Sprintf("%s/%s/runs", generator.ThreadUrl, thread.Id)

	req, err := http.NewRequest("POST", runUrl, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return nil, errors.New("there was a problem creating the request: " + err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+generator.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New("there was a problem making the request: " + err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("there was a problem making the request: " + response.Status)
	}

	runBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("there was a problem reading body: " + err.Error())
	}

	var runData Run
	err = json.Unmarshal(runBody, &runData)
	if err != nil {
		return nil, errors.New("there was a problem parsing the json: " + err.Error())
	}
	return &runData, nil
}

func (run *Run) pollAndCheckRun(httpClient http.Client, threadUrl, apiKey string) error {
	viableStatuses := []string{"completed", "cancelled", "failed", "incomplete", "expired"}

	for !slices.Contains(viableStatuses, run.Status) {
		time.Sleep(5 * time.Second)
		newRun, err := run.checkRunStatus(httpClient, threadUrl, apiKey)
		if err != nil {
			return err
		}
		*run = *newRun
	}

	if run.Status != "completed" {
		return fmt.Errorf("run ended with status: %s", run.Status)
	}
	return nil
}

func (run *Run) checkRunStatus(httpClient http.Client, threadUrl, apiKey string) (*Run, error) {
	url := fmt.Sprintf("%s/%s/runs/%s", threadUrl, run.ThreadId, run.Id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var newRun Run
	err = json.Unmarshal(body, &newRun)
	if err != nil {
		return nil, err
	}

	return &newRun, nil
}

func (thread *Thread) listMessages(httpClient http.Client, threadUrl, apiKey string) (*ListMessagesResponse, error) {
	listUrl := fmt.Sprintf("%s/%s/messages?order=desc", threadUrl, thread.Id)
	req, err := http.NewRequest("GET", listUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("there was a problem making the request: " + response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var messageObject ListMessagesResponse
	err = json.Unmarshal(body, &messageObject)
	if err != nil {
		return nil, err
	}
	return &messageObject, nil
}

func (messageObject *ListMessagesResponse) getFinalMessageContent() (*PromptRecipe, error) {
	var recipe PromptRecipe
	err := json.Unmarshal([]byte(messageObject.Data[0].Content[0].Text.Value), &recipe)
	if err != nil {
		return nil, errors.New("there was a problem parsing the json: " + err.Error())
	}

	return &recipe, nil
}

func (generator *OpenAIPromptGenerator) GenerateRecipe(prompt string) (*PromptCocktail, error) {
	/*
		Steps:
		1. Create a new thread
		2. Create a new message to the thread
		3. Wait for the response of the assistant
		4. Return the prompt recipe
	*/
	httpClient := &http.Client{}

	thread, err := generator.createThread(*httpClient)
	if err != nil {
		return nil, errors.New("there was a problem creating a thread")
	}

	err = thread.addMessageToThread(*httpClient, prompt, generator.ThreadUrl, generator.ApiKey)
	if err != nil {
		return nil, errors.New("there was a problem adding a message to the thread: " + err.Error())
	}

	run, err := generator.createRun(*httpClient, thread)
	if err != nil {
		return nil, err
	}

	err = run.pollAndCheckRun(*httpClient, generator.ThreadUrl, generator.ApiKey)
	if err != nil {
		return nil, errors.New("there was a problem with the current run: " + err.Error())
	}

	messages, err := thread.listMessages(*httpClient, generator.ThreadUrl, generator.ApiKey)
	if err != nil {
		return nil, errors.New("there was a problem listing messages")
	}

	promptRecipe, err := messages.getFinalMessageContent()
	if err != nil {
		return nil, errors.New("there was a problem getting the final message content, " + err.Error())
	}

	return &promptRecipe.Cocktail, nil
}
