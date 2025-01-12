package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-chat/shared"
	"io"
	"log"
	"net/http"
)

// makeRequest makes a type-safe HTTP request to the given route
func makeRequest[Req any, Res any](route shared.Route[Req, Res], req Req) (*Res, error) {
	var body io.Reader
	if !isEmptyStruct(req) {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	httpReq, err := http.NewRequest(route.Method, route.Path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Res
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// isEmptyStruct returns true if the given value is an empty struct
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

// FetchHealthStatus fetches the current health status from the server
func FetchHealthStatus() (*shared.HealthResponse, error) {
	health, err := makeRequest(shared.RouteHealth, struct{}{})
	if err != nil {
		log.Printf("❌ Error fetching health status: %v", err)
		return nil, err
	}

	log.Printf("✅ Health status: %+v", health)
	return health, nil
}
