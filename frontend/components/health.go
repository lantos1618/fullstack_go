package components

import (
	"go-chat/frontend/store/actions"
	"go-chat/shared"
	"log"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// HealthStatus is a component that displays the server health status
type HealthStatus struct {
	vecty.Core
	health *shared.HealthResponse
}

// Mount implements the vecty.Mounter interface
func (h *HealthStatus) Mount() {
	go h.fetchHealth()
}

func (h *HealthStatus) fetchHealth() {
	health, err := actions.FetchHealthStatus()
	if err != nil {
		log.Printf("❌ Failed to fetch health status: %v", err)
		return
	}
	h.health = health
	vecty.Rerender(h)
}

// Render implements the vecty.Component interface
func (h *HealthStatus) Render() vecty.ComponentOrHTML {
	var status string
	var timestamp string
	var statusClass string

	if h.health == nil {
		status = "Loading..."
		timestamp = "N/A"
		statusClass = "text-gray-500"
	} else {
		status = h.health.Status
		timestamp = time.Unix(h.health.Timestamp, 0).Format(time.RFC3339)
		if status == "ok" {
			statusClass = "text-green-500"
		} else {
			statusClass = "text-red-500"
		}
	}

	return elem.Div(
		vecty.Markup(
			vecty.Class("flex", "items-center", "space-x-2", "text-sm"),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("font-medium"),
			),
			vecty.Text("Server Status:"),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class(statusClass, "font-bold"),
			),
			vecty.Text(status),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("text-gray-500", "dark:text-gray-400"),
			),
			vecty.Text(timestamp),
		),
	)
}
