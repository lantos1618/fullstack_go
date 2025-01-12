package ws

// MessageType represents different types of WebSocket messages
type MessageType string

const (
	// System message types
	TypePing  MessageType = "PING"  // Ping message to check connection
	TypePong  MessageType = "PONG"  // Pong response to ping
	TypeClose MessageType = "CLOSE" // Connection close message
	TypeError MessageType = "ERROR" // Error message

	// Chat message types
	TypeMessage MessageType = "MESSAGE" // Regular chat message
	TypeTyping  MessageType = "TYPING"  // Typing indicator
	TypeJoin    MessageType = "JOIN"    // User joined notification
	TypeLeave   MessageType = "LEAVE"   // User left notification
)

// Message represents a WebSocket message structure
type Message struct {
	Type    MessageType `json:"type"`              // Type of the message
	Payload any         `json:"payload,omitempty"` // Message payload, can be any of the payload types below
}

// Payload types for different message types

// TextMessage is the payload for TypeMessage
type TextMessage struct {
	Text string `json:"text"` // The actual message text
	From string `json:"from"` // Username of the sender
}

// ErrorMessage is the payload for TypeError
type ErrorMessage struct {
	Error string `json:"error"` // Error description
}

// TypingMessage is the payload for TypeTyping
type TypingMessage struct {
	From     string `json:"from"`      // Username of the person typing
	IsTyping bool   `json:"is_typing"` // Whether the user is typing
}

// JoinMessage is the payload for TypeJoin
type JoinMessage struct {
	From string `json:"from"` // Username of the person who joined
}

// LeaveMessage is the payload for TypeLeave
type LeaveMessage struct {
	From string `json:"from"` // Username of the person who left
}
