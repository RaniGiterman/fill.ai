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
