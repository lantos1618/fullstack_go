//go:build wasm
// +build wasm

package actions

import (
	"go-chat/shared/api"
	"go-chat/shared/http"
	"log"
)

// healthClient is a type-safe client for health check requests
var healthClient = http.NewClient[struct{}, api.HealthResponse]("", api.RouteHealth)

// FetchHealthStatus fetches the current health status from the server
func FetchHealthStatus() (*api.HealthResponse, error) {
	health, err := healthClient.Request(struct{}{})
	if err != nil {
		log.Printf("❌ Error fetching health status: %v", err)
		return nil, err
	}

	log.Printf("✅ Health status: %+v", health)
	return health, nil
}
