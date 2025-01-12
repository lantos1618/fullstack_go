package api

import (
	"go-chat/shared/http"
)

// WebSocket message types
const (
	WSTypeMessage = "message"
	WSTypeTyping  = "typing"
)

// Common response types
type (
	// ErrorResponse represents a standardized error response
	ErrorResponse struct {
		Error   string `json:"error"`
		Code    string `json:"code"`
		Details string `json:"details,omitempty"`
	}

	// HealthResponse represents the health check response
	HealthResponse struct {
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
	}

	// ChatMessage represents a chat message
	ChatMessage struct {
		ID        string `json:"id"`
		UserID    string `json:"userId"`
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp"`
	}

	// SendMessageRequest represents the request to send a message
	SendMessageRequest struct {
		Content string `json:"content"`
	}

	// WSMessage represents a WebSocket message envelope
	WSMessage struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}

	// WSChatMessage represents a chat message sent over WebSocket
	WSChatMessage struct {
		Text string `json:"text"`
		From string `json:"from"`
	}

	// WSTypingMessage represents a typing indicator message
	WSTypingMessage struct {
		From     string `json:"from"`
		IsTyping bool   `json:"isTyping"`
	}
)

// API Routes - Single source of truth for all API endpoints
var (
	// Health Routes
	RouteHealth = http.NewRoute[struct{}, HealthResponse]("/api/health", http.MethodGet)

	// Chat Routes
	RouteSendMessage   = http.NewRoute[SendMessageRequest, ChatMessage]("/api/chat/messages", http.MethodPost)
	RouteGetMessages   = http.NewRoute[struct{}, []ChatMessage]("/api/chat/messages", http.MethodGet)
	RouteWebSocketChat = http.NewRoute[struct{}, struct{}]("/api/chat/ws", http.MethodGet)
)

// APIError codes for standardized error handling
const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeServerError    = "SERVER_ERROR"
)
