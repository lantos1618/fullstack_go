package shared

// MessageType represents different types of WebSocket messages
type MessageType string

const (
	// Message types
	TypePing    MessageType = "PING"
	TypePong    MessageType = "PONG"
	TypeMessage MessageType = "MESSAGE"
	TypeError   MessageType = "ERROR"
)

// WSMessage represents a WebSocket message structure
type WSMessage struct {
	Type    MessageType `json:"type"`
	Payload any         `json:"payload"`
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
