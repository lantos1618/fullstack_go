//go:build wasm
// +build wasm

package components

import (
	"encoding/json"
	"log"
	"syscall/js"

	"go-chat/frontend/actions"
	"go-chat/frontend/dispatcher"
	"go-chat/frontend/store"
	"go-chat/shared"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
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
	log.Printf("Chat component mounted")
	store.Listeners.Add(c, func() {
		log.Printf("Chat component rerendering due to store change")
		vecty.Rerender(c)
	})
}

// Unmount implements the vecty.Unmounter interface
func (c *Chat) Unmount() {
	log.Printf("Chat component unmounted")
	store.Listeners.Remove(c)
}

var chatInstance *Chat

// NewChat creates a new chat component
func NewChat() vecty.ComponentOrHTML {
	log.Printf("Creating new chat component")
	if chatInstance == nil {
		chatInstance = &Chat{}
		chatInstance.connectWS()
	}
	return chatInstance
}

func (c *Chat) connectWS() {
	ws := js.Global().Get("WebSocket").New("ws://" + js.Global().Get("location").Get("host").String() + "/ws")

	ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data").String()
		var wsMsg shared.WSMessage
		if err := json.Unmarshal([]byte(data), &wsMsg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			return nil
		}

		switch wsMsg.Type {
		case shared.TypeMessage:
			if textMsg, ok := wsMsg.Payload.(map[string]interface{}); ok {
				dispatcher.Dispatch(&actions.AddMessage{
					Text: textMsg["text"].(string),
					From: textMsg["from"].(string),
				})
			}
		case shared.TypeTyping:
			var typingMsg shared.TypingMessage
			if payloadBytes, err := json.Marshal(wsMsg.Payload); err == nil {
				if err := json.Unmarshal(payloadBytes, &typingMsg); err == nil {
					if typingMsg.From != store.Username {
						dispatcher.Dispatch(&actions.SetTyping{
							Username: typingMsg.From,
							IsTyping: true,
						})
						// Clear typing indicator after 3 seconds
						js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
							dispatcher.Dispatch(&actions.SetTyping{
								Username: typingMsg.From,
								IsTyping: false,
							})
							return nil
						}), 3000)
					}
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
	msg := shared.WSMessage{
		Type: shared.TypeTyping,
		Payload: shared.TypingMessage{
			From: store.Username,
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

	msg := shared.WSMessage{
		Type: shared.TypeMessage,
		Payload: shared.TextMessage{
			Text: c.input,
			From: store.Username,
		},
	}

	if data, err := json.Marshal(msg); err == nil {
		c.ws.Call("send", string(data))
	}

	c.input = ""
	vecty.Rerender(c)
}

func (c *Chat) renderMessageList() vecty.ComponentOrHTML {
	if len(store.Messages) == 0 {
		return elem.Paragraph(
			vecty.Markup(
				vecty.Class("text-gray-500", "text-center"),
			),
			vecty.Text("No messages yet"),
		)
	}

	var messageElements []vecty.MarkupOrChild
	for _, msg := range store.Messages {
		messageElements = append(messageElements,
			elem.Div(
				vecty.Markup(
					vecty.Class("mb-2"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("font-bold"),
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
			vecty.Class("text-gray-500", "italic", "text-sm", "mt-2"),
		),
		vecty.Text(typingText),
	)
}

// Render implements the vecty.Component interface
func (c *Chat) Render() vecty.ComponentOrHTML {
	log.Printf("Chat.Render called")
	result := elem.Div(
		vecty.Markup(
			vecty.Class("container", "mx-auto", "p-4"),
		),
		elem.Heading1(
			vecty.Markup(
				vecty.Class("text-2xl", "font-bold", "mb-4"),
			),
			vecty.Text("Chat Application"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("bg-white", "rounded-lg", "shadow-md", "p-4", "mb-4", "h-96", "overflow-y-auto"),
			),
			c.renderMessageList(),
			c.renderTypingIndicators(),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("flex", "gap-2"),
			),
			elem.Input(
				vecty.Markup(
					vecty.Class("flex-1", "p-2", "border", "rounded"),
					event.Input(c.onInput),
					vecty.Property("type", "text"),
					vecty.Property("value", c.input),
					vecty.Property("placeholder", "Type a message..."),
				),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class("bg-blue-500", "text-white", "px-4", "py-2", "rounded", "hover:bg-blue-600"),
					event.Click(c.onSend),
				),
				vecty.Text("Send"),
			),
		),
	)
	log.Printf("Chat.Render completed")
	return result
}
