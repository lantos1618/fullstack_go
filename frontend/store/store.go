//go:build wasm
// +build wasm

package store

import (
	"go-chat/frontend/actions"
	"go-chat/frontend/dispatcher"
	"go-chat/shared"
	"log"
	"syscall/js"
)

var (
	// Messages represents all chat messages
	Messages []shared.TextMessage

	// Username represents the current user's username
	Username string

	// TypingUsers represents users who are currently typing
	TypingUsers = make(map[string]bool)

	// IsDarkMode represents the current theme state
	IsDarkMode bool

	// Listeners will be invoked when the store changes
	Listeners = NewListenerRegistry()
)

func init() {
	// Load dark mode preference from localStorage
	localStorage := js.Global().Get("localStorage")
	darkMode := localStorage.Call("getItem", "darkMode")
	if !darkMode.IsNull() && darkMode.String() == "true" {
		IsDarkMode = true
		log.Printf("💾 Loaded dark mode from storage: %v", IsDarkMode)
	}

	log.Printf("💾 Store initialized | darkMode: %v", IsDarkMode)
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
		log.Printf("👤 Username set to: %s", Username)

	case *actions.AddMessage:
		Messages = append(Messages, shared.TextMessage{
			Text: a.Text,
			From: a.From,
		})
		log.Printf("💬 Message added from %s: %s", a.From, a.Text)

	case *actions.SetTyping:
		if a.IsTyping {
			TypingUsers[a.Username] = true
			log.Printf("⌨️  %s started typing", a.Username)
		} else {
			delete(TypingUsers, a.Username)
			log.Printf("⌨️  %s stopped typing", a.Username)
		}

	case *actions.ToggleDarkMode:
		IsDarkMode = !IsDarkMode
		// Save dark mode preference to localStorage
		localStorage := js.Global().Get("localStorage")
		localStorage.Call("setItem", "darkMode", js.ValueOf(IsDarkMode))
		log.Printf("🎨 Dark mode toggled | new state: %v | saved to storage", IsDarkMode)

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}
