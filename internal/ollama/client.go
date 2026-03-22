package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is an HTTP client for the Ollama API.
type Client struct {
	baseURL string
	model   string
	http    *http.Client
}

// NewClient creates a new Ollama client.
func NewClient(baseURL, model string) *Client {
	return &Client{
		baseURL: baseURL,
		model:   model,
		http: &http.Client{
			Timeout: 600 * time.Second,
		},
	}
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
}

// Generate sends a prompt to Ollama and returns the response text.
func (c *Client) Generate(system, prompt string) (string, error) {
	req := generateRequest{
		Model:  c.model,
		Prompt: prompt,
		System: system,
		Stream: false,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	resp, err := c.http.Post(c.baseURL+"/api/generate", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("ollama: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: returned %d", resp.StatusCode)
	}

	var result generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ollama: decode response: %w", err)
	}
	return result.Response, nil
}

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type chatResponse struct {
	Message Message `json:"message"`
}

// Chat sends a conversation history to Ollama and returns the assistant's response.
// The context controls cancellation — if ctx is cancelled the HTTP request is aborted.
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	req := chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("ollama: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ollama: chat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: chat returned %d", resp.StatusCode)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ollama: decode chat response: %w", err)
	}
	return result.Message.Content, nil
}
