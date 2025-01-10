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
