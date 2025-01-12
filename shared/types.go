package shared

// WsMessageType represents different types of WebSocket messages
type WsMessageType string

const (
	// Message types
	TypePing    WsMessageType = "PING"
	TypePong    WsMessageType = "PONG"
	TypeMessage WsMessageType = "MESSAGE"
	TypeError   WsMessageType = "ERROR"
	TypeClose   WsMessageType = "CLOSE"
	TypeJoin    WsMessageType = "JOIN"
	TypeLeave   WsMessageType = "LEAVE"
	TypeTyping  WsMessageType = "TYPING"
)

// Route represents an API route with its associated request and response types
type Route[Req any, Res any] struct {
	Path   string
	Method string
}

// NewRoute creates a new route with the given path and method
func NewRoute[Req any, Res any](path, method string) Route[Req, Res] {
	return Route[Req, Res]{
		Path:   path,
		Method: method,
	}
}

// API Routes
var (
	// RouteHealth represents the health check endpoint
	// It takes no request body (struct{}) and returns HealthResponse
	RouteHealth = NewRoute[struct{}, HealthResponse]("/api/health", "GET")
)

// WSMessage represents a WebSocket message structure
type WSMessage struct {
	Type    WsMessageType `json:"type"`
	Payload any           `json:"payload"`
}

// TextMessage represents a simple text message
type TextMessage struct {
	Text string `json:"text"`
	From string `json:"from"`
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	Error string `json:"error"`
}

// TypingMessage represents a typing message
type TypingMessage struct {
	From string `json:"from"`
}

// JoinMessage represents a join message
type JoinMessage struct {
	From string `json:"from"`
}

// LeaveMessage represents a leave message
type LeaveMessage struct {
	From string `json:"from"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}
