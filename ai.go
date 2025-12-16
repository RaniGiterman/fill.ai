package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

const SYSTEM_PROMPT = `
You are a precise data extraction assistant. You will receive an HTML page and must extract information to fill a JSON object with exactly three fields.

## Input Format
You will receive:
1. An empty JSON template: {"title": "", "description": "", "img": ""}
2. A complete HTML page

## Field Descriptions

**title** (string):
- Extract the main page title
- Check in order: <title> tag, <h1> tag, og:title meta tag, or the most prominent heading
- Return only the text content, no HTML tags
- If multiple candidates exist, choose the most descriptive and prominent one

**description** (string):
- Extract a brief summary or description of the page content
- Check in order: meta description tag, og:description, first <p> tag with substantial content, or summarize the main content in 1-2 sentences
- Maximum ~150-200 characters preferred
- Should describe what the page is about

**img** (string):
- Extract the primary/featured image URL
- Check in order: og:image meta tag, twitter:image, first prominent <img> src in main content area, or hero image
- Return the complete URL (absolute path)
- If relative URL, construct absolute URL using the page's domain
- Return empty string "" if no suitable image found

## Output Requirements
- Return ONLY the filled JSON object
- No explanation, no markdown, no additional text
- Ensure valid JSON syntax
- Use empty string "" for any field that cannot be determined

## Example Output Format
{"title": "Example Page Title", "description": "A brief description of what this page contains and its main purpose.", "img": "https://example.com/images/featured.jpg"}

Now process the provided HTML and return the filled JSON.

`

func queryLLM(jsonInput, html string) (string, error) {
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
