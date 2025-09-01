package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"nextjs-to-openapi/internal/models"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	model      string
}

func NewClient(baseURL, model string) *Client {
	return &Client{
		baseURL:    baseURL,
		model:      model,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// RouteDocumentation represents the structured response we want from Ollama
type RouteDocumentation struct {
	Path        string            `json:"path"`
	Methods     map[string]Method `json:"methods"`
	Description string            `json:"description"`
}

// Method represents an HTTP method documentation
type Method struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	In       string `json:"in"` // "path", "query", "body"
	Required bool   `json:"required"`
}

// DocumentRoute sends a route to Ollama for documentation
func (c *Client) DocumentRoute(route models.APIRoute) (*RouteDocumentation, error) {
	prompt := c.buildPrompt(route)

	// Send request to Ollama
	response, err := c.sendRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}

	// Parse the response
	doc, err := c.parseResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return doc, nil
}

// buildPrompt creates a smart prompt for Ollama
func (c *Client) buildPrompt(route models.APIRoute) string {
	return fmt.Sprintf(`Analyze this Next.js API route file and extract OpenAPI information.

File: %s
File Type: %s
Content:
%s

IMPORTANT: Return ONLY valid JSON with no markdown formatting, no backticks, no code blocks.

Return this exact JSON structure:
{
  "path": "/api/path/here",
  "description": "Brief description of what this API endpoint does",
  "methods": {
    "GET": {
      "summary": "Brief summary",
      "description": "Detailed description", 
      "parameters": [
        {
          "name": "paramName",
          "type": "string",
          "in": "path",
          "required": true
        }
      ]
    }
  }
}

Rules:
1. Convert [id] to {id} in the path
2. Convert [...slug] to {slug} in the path
3. Only include methods that actually exist in the code
4. Return ONLY the JSON, no markdown, no explanations, no code blocks
`, route.FilePath, route.FileType, route.Content)
}

// sendRequest sends the prompt to Ollama
func (c *Client) sendRequest(prompt string) (string, error) {
	// Create request payload
	reqPayload := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false, // We want the complete response at once
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.baseURL + "/api/generate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

// parseResponse attempts to extract JSON from Ollama's response
func (c *Client) parseResponse(response string) (*RouteDocumentation, error) {
	fmt.Printf("\n🐛 DEBUG - Raw Ollama response:\n")
	fmt.Printf("--- START RESPONSE ---\n")
	fmt.Printf("%s\n", response)
	fmt.Printf("--- END RESPONSE ---\n\n")

	// Clean up the response - remove markdown code blocks
	cleanedResponse := cleanMarkdownJSON(response)

	fmt.Printf("🐛 DEBUG - Cleaned response:\n")
	fmt.Printf("--- START CLEANED ---\n")
	fmt.Printf("%s\n", cleanedResponse)
	fmt.Printf("--- END CLEANED ---\n\n")

	var doc RouteDocumentation
	if err := json.Unmarshal([]byte(cleanedResponse), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &doc, nil
}

// cleanMarkdownJSON removes markdown code block formatting from JSON
func cleanMarkdownJSON(response string) string {
	// Remove various markdown patterns

	response = strings.ReplaceAll(response, "```JSON", "")

	// Remove any leading/trailing whitespace
	response = strings.TrimSpace(response)

	// Find the first { and last } to extract just the JSON part
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		response = response[start : end+1]
	}

	return response
}
