package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s - %s", e.StatusCode, e.Message, e.Details)
}

// Client provides type-safe HTTP request methods
type Client[Req any, Res any] struct {
	baseURL    string
	httpClient *http.Client
	route      Route[Req, Res]
}

// NewClient creates a new HTTP client with the given base URL and route
func NewClient[Req any, Res any](baseURL string, route Route[Req, Res]) *Client[Req, Res] {
	return &Client[Req, Res]{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		route:      route,
	}
}

// Request makes a type-safe HTTP request using the client's route
func (c *Client[Req, Res]) Request(req Req) (*Res, error) {
	var body io.Reader
	if !isEmptyStruct(req) {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + c.route.Path
	httpReq, err := http.NewRequest(c.route.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		if err := json.Unmarshal(respBody, &apiErr); err == nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    apiErr.Error,
				Details:    apiErr.Details,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
			Details:    string(respBody),
		}
	}

	// Parse successful response
	var result Res
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// isEmptyStruct returns true if the value is an empty struct
func isEmptyStruct(v any) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case struct{}:
		return true
	default:
		return false
	}
}
