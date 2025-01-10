package store

import (
	"go-chat/frontend/actions"
	"go-chat/frontend/dispatcher"
	"go-chat/shared"
)

var (
	// Messages represents all chat messages
	Messages []shared.TextMessage

	// Username represents the current user's username
	Username string

	// TypingUsers represents users who are currently typing
	TypingUsers = make(map[string]bool)

	// Listeners will be invoked when the store changes
	Listeners = NewListenerRegistry()
)

func init() {
	dispatcher.Register(onAction)
}

// NewListenerRegistry creates a new listener registry
func NewListenerRegistry() *ListenerRegistry {
	return &ListenerRegistry{
		listeners: make(map[interface{}]func()),
	}
}

// ListenerRegistry manages store listeners
type ListenerRegistry struct {
	listeners map[interface{}]func()
}

// Add adds a listener with a key
func (r *ListenerRegistry) Add(key interface{}, listener func()) {
	if key == nil {
		key = new(int)
	}
	r.listeners[key] = listener
}

// Remove removes a listener by key
func (r *ListenerRegistry) Remove(key interface{}) {
	delete(r.listeners, key)
}

// Fire invokes all listeners
func (r *ListenerRegistry) Fire() {
	for _, l := range r.listeners {
		l()
	}
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.SetUsername:
		Username = a.Username

	case *actions.AddMessage:
		Messages = append(Messages, shared.TextMessage{
			Text: a.Text,
			From: a.From,
		})

	case *actions.SetTyping:
		if a.IsTyping {
			TypingUsers[a.Username] = true
		} else {
			delete(TypingUsers, a.Username)
		}

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}
