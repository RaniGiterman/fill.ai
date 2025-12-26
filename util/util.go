package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
)

const SYSTEM_PROMPT = `
You are a precise data extraction assistant. You will receive an HTML page that describes a consumer product with details and must extract information to fill a JSON object with exactly three fields.

## Input Format
You will receive:
1. An empty JSON template: {"title": "", "description": "", "img": ""}
2. A base URL for the website
3. A complete HTML page

## Field Descriptions

*title* (string):
- Extract the product title
- Return only the text content, no HTML tags
- If multiple candidates exist, choose the most descriptive and prominent one

*description* (string):
- Extract product description from the page content

*img* (string):
- Extract the product image. if more than one, then concat them all with comma
- If image URL is relative (starts with / or does not include http/https), prepend the base URL to create absolute URL
- Example: if base URL is "https://example.com" and img is "/images/product.png", return "https://example.com/images/product.png"

## Output Requirements
- Return ONLY the filled JSON object
- Ensure valid JSON syntax
- Use empty string "" for any field that cannot be determined

Now process the provided HTML and return the filled JSON.
`

// Run HTML against openai model
func QueryGPT(ctx context.Context, client openai.Client, html, url string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(SYSTEM_PROMPT),
			openai.UserMessage(url),
			openai.UserMessage(html),
		},
		Seed:            openai.Int(0),
		Model:           openai.ChatModelGPT5Mini,
		ReasoningEffort: openai.ReasoningEffortLow,
	}

	// completion, err := client.Chat.Completions.New(ctx, params)
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	return completion.Choices[0].Message.Content, nil
}

// Given URL, return page full HTML
func QueryHTML(url string) (string, error) {
	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	// Ensure the response body is closed after the function returns
	defer resp.Body.Close()

	// Check if the request was successful (HTTP status code 200)
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print the page content (the HTML body)
	return string(body), nil
}

// Run against Ollama model, not used
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func queryOllama(jsonInput, html string) (string, error) {
	userPrompt := fmt.Sprintf(`Empty JSON template:
	%s

HTML page:
%s`, jsonInput, html)

	reqBody := OllamaRequest{
		Model:  "model460",
		Prompt: userPrompt,
		System: SYSTEM_PROMPT,
		Stream: false,
	}

	// turn request body to json from golang struct
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Minute, // Increase timeout for LLM processing
	}

	// Send the request to Ollama API
	resp, err := client.Post("http://192.168.1.3:6767/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return ollamaResp.Response, nil
}
