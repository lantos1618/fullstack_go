//go:build wasm
// +build wasm

package components

import (
	"encoding/json"
	"log"
	"syscall/js"

	"go-chat/frontend/store"
	"go-chat/frontend/store/actions"
	"go-chat/frontend/store/dispatcher"
	"go-chat/shared/ws"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
)

// Chat is the main chat component
type Chat struct {
	vecty.Core
	ws          js.Value
	input       string
	typingTimer js.Value
}

// Mount implements the vecty.Mounter interface
func (c *Chat) Mount() {
	log.Printf("🚀 Chat component mounted")
	store.Listeners.Add(c, func() {
		log.Printf("🔄 Chat component rerendering due to store change")
		vecty.Rerender(c)
	})
}

// Unmount implements the vecty.Unmounter interface
func (c *Chat) Unmount() {
	log.Printf("👋 Chat component unmounted")
	store.Listeners.Remove(c)
	chatInstance = nil // Reset the singleton instance
}

var chatInstance *Chat

// NewChat creates a new chat component
func NewChat() vecty.ComponentOrHTML {
	if chatInstance == nil {
		log.Printf("📱 Creating new chat component")
		chatInstance = &Chat{}
		chatInstance.connectWS()
	}
	return chatInstance
}

func (c *Chat) connectWS() {
	ws := js.Global().Get("WebSocket").New("ws://" + js.Global().Get("location").Get("host").String() + "/ws")

	ws.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("WebSocket connection established")
		return nil
	}))

	ws.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("WebSocket error occurred")
		return nil
	}))

	ws.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("WebSocket connection closed, attempting to reconnect in 3 seconds...")
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.connectWS()
			return nil
		}), 3000)
		return nil
	}))

	ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data").String()
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			return nil
		}

		msgType, ok := msg["type"].(string)
		if !ok {
			log.Printf("invalid message type")
			return nil
		}

		log.Printf("Received message of type: %s", msgType)

		switch msgType {
		case "MESSAGE":
			if payload, ok := msg["payload"].(map[string]interface{}); ok {
				text, _ := payload["text"].(string)
				from, _ := payload["from"].(string)
				if text != "" && from != "" {
					log.Printf("Received chat message from %s: %s", from, text)
					dispatcher.Dispatch(&actions.AddMessage{
						Text: text,
						From: from,
					})
				}
			}
		case "TYPING":
			if payload, ok := msg["payload"].(map[string]interface{}); ok {
				from, _ := payload["from"].(string)
				isTyping, _ := payload["is_typing"].(bool)
				if from != "" && from != store.Username {
					dispatcher.Dispatch(&actions.SetTyping{
						Username: from,
						IsTyping: isTyping,
					})
					// Clear typing indicator after 1 second
					js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						dispatcher.Dispatch(&actions.SetTyping{
							Username: from,
							IsTyping: false,
						})
						return nil
					}), 1000)
				}
			}
		}
		return nil
	}))

	c.ws = ws
}

func (c *Chat) onInput(e *vecty.Event) {
	c.input = e.Target.Get("value").String()
	vecty.Rerender(c)

	// Send typing notification
	msg := map[string]interface{}{
		"type": "TYPING",
		"payload": map[string]interface{}{
			"from":      store.Username,
			"is_typing": true,
		},
	}
	if data, err := json.Marshal(msg); err == nil {
		c.ws.Call("send", string(data))
	}
}

func (c *Chat) onSend(e *vecty.Event) {
	if c.input == "" {
		return
	}

	msg := map[string]interface{}{
		"type": "MESSAGE",
		"payload": map[string]interface{}{
			"text": c.input,
			"from": store.Username,
		},
	}

	if data, err := json.Marshal(msg); err == nil {
		c.ws.Call("send", string(data))
		log.Printf("Sent message: %s", c.input)
	}

	c.input = ""
	vecty.Rerender(c)
}

func (c *Chat) onKeyDown(e *vecty.Event) {
	// Check if the pressed key is Enter
	if e.Get("key").String() == "Enter" {
		// Check if Shift is held down
		if e.Get("shiftKey").Bool() {
			// Add newline character
			c.input += "\n"
			vecty.Rerender(c)
		} else {
			// Prevent default behavior (newline)
			e.Call("preventDefault")
			// Submit the message
			c.onSend(e)
		}
	}
}

func (c *Chat) renderMessageList() vecty.ComponentOrHTML {
	if len(store.Messages) == 0 {
		return elem.Paragraph(
			vecty.Markup(
				vecty.Class("text-gray-500", "dark:text-gray-400", "text-center", "italic"),
			),
			vecty.Text("No messages yet"),
		)
	}

	var messageElements []vecty.MarkupOrChild
	for _, msg := range store.Messages {
		messageElements = append(messageElements,
			elem.Div(
				vecty.Markup(
					vecty.Class("mb-4", "text-gray-800", "dark:text-gray-200"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("font-bold", "text-blue-600", "dark:text-blue-400", "mr-2"),
					),
					vecty.Text(msg.From+": "),
				),
				vecty.Text(msg.Text),
			),
		)
	}
	return elem.Div(messageElements...)
}

func (c *Chat) renderTypingIndicators() vecty.ComponentOrHTML {
	if len(store.TypingUsers) == 0 {
		return nil
	}

	var typingText string
	var users []string
	for user := range store.TypingUsers {
		users = append(users, user)
	}

	if len(users) == 1 {
		typingText = users[0] + " is typing..."
	} else {
		typingText = "Multiple people are typing..."
	}

	return elem.Paragraph(
		vecty.Markup(
			vecty.Class("text-gray-500", "dark:text-gray-400", "italic", "text-sm", "mt-2"),
		),
		vecty.Text(typingText),
	)
}

// Render implements the vecty.Component interface
func (c *Chat) Render() vecty.ComponentOrHTML {
	log.Printf("🎨 Chat component rendering")
	result := elem.Div(
		vecty.Markup(
			vecty.Class("container", "mx-auto", "p-4", "max-w-4xl"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class(
					"bg-white", "dark:bg-gray-800",
					"rounded-lg", "shadow-lg",
					"p-6", "mb-4", "h-96", "overflow-y-auto",
					"transition-colors", "duration-200",
					"border", "border-gray-200", "dark:border-gray-700",
				),
			),
			c.renderMessageList(),
			c.renderTypingIndicators(),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("flex", "gap-3"),
			),
			elem.Input(
				vecty.Markup(
					vecty.Class(
						"flex-1", "p-3",
						"border", "border-gray-300", "dark:border-gray-600",
						"rounded-lg",
						"bg-white", "dark:bg-gray-700",
						"text-gray-900", "dark:text-white",
						"placeholder-gray-500", "dark:placeholder-gray-400",
						"focus:ring-2", "focus:ring-blue-500", "dark:focus:ring-blue-400",
						"focus:border-transparent",
						"transition-colors", "duration-200",
					),
					event.Input(c.onInput),
					event.KeyDown(c.onKeyDown),
					prop.Value(c.input),
					prop.Placeholder("Type a message..."),
				),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class(
						"px-6", "py-3",
						"bg-blue-500", "dark:bg-blue-600",
						"text-white",
						"rounded-lg",
						"hover:bg-blue-600", "dark:hover:bg-blue-700",
						"focus:outline-none", "focus:ring-2", "focus:ring-blue-500",
						"transition-colors", "duration-200",
						"disabled:opacity-50",
					),
					event.Click(c.onSend),
					prop.Disabled(c.input == ""),
				),
				vecty.Text("Send"),
			),
		),
	)

	return result
}

// SendMessage sends a chat message through the WebSocket connection
func (c *Chat) SendMessage(text string) {
	if c.ws.Get("readyState").Int() != 1 { // 1 = OPEN
		log.Printf("Cannot send message: WebSocket is not connected")
		return
	}

	msg := ws.Message{
		Type: ws.TypeMessage,
		Payload: ws.TextMessage{
			Text: text,
			From: store.Username,
		},
	}

	if msgBytes, err := json.Marshal(msg); err == nil {
		c.ws.Call("send", string(msgBytes))
		log.Printf("Sent message: %s", text)
	} else {
		log.Printf("error marshaling message: %v", err)
	}
}
