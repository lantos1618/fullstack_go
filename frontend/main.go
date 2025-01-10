//go:build wasm
// +build wasm

package main

import (
	"encoding/json"
	"log"
	"syscall/js"

	"go-chat/shared"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type App struct {
	vecty.Core
	ws       js.Value
	messages []shared.TextMessage
	input    string
}

func NewApp() *App {
	app := &App{}
	app.connectWS()
	return app
}

func (a *App) connectWS() {
	ws := js.Global().Get("WebSocket").New("ws://" + js.Global().Get("location").Get("host").String() + "/ws")

	ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data").String()
		var wsMsg shared.WSMessage
		if err := json.Unmarshal([]byte(data), &wsMsg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			return nil
		}

		if wsMsg.Type == shared.TypeMessage {
			if textMsg, ok := wsMsg.Payload.(map[string]interface{}); ok {
				a.messages = append(a.messages, shared.TextMessage{
					Text: textMsg["text"].(string),
					From: textMsg["from"].(string),
				})
				vecty.Rerender(a)
			}
		}
		return nil
	}))

	a.ws = ws
}

func (a *App) sendMessage(e *vecty.Event) {
	if a.input == "" {
		return
	}

	msg := shared.WSMessage{
		Type: shared.TypeMessage,
		Payload: shared.TextMessage{
			Text: a.input,
			From: "User",
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	a.ws.Call("send", string(data))
	a.input = ""
	vecty.Rerender(a)
}

func (a *App) onInput(e *vecty.Event) {
	a.input = e.Target.Get("value").String()
	vecty.Rerender(a)
}

func (a *App) renderMessageList() vecty.ComponentOrHTML {
	if len(a.messages) == 0 {
		return elem.Paragraph(
			vecty.Markup(
				vecty.Class("text-gray-500 text-center"),
			),
			vecty.Text("No messages yet"),
		)
	}

	var messageElements []vecty.MarkupOrChild
	for _, msg := range a.messages {
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

func (a *App) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(
			vecty.Class("min-h-screen bg-gray-100"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("container mx-auto p-4"),
			),
			elem.Heading1(
				vecty.Markup(
					vecty.Class("text-2xl font-bold mb-4"),
				),
				vecty.Text("Chat Application"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("bg-white rounded-lg shadow-md p-4 mb-4 h-96 overflow-y-auto"),
				),
				a.renderMessageList(),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("flex gap-2"),
				),
				elem.Input(
					vecty.Markup(
						vecty.Class("flex-1 p-2 border rounded"),
						event.Input(a.onInput),
						vecty.Property("type", "text"),
						vecty.Property("value", a.input),
						vecty.Property("placeholder", "Type a message..."),
					),
				),
				elem.Button(
					vecty.Markup(
						vecty.Class("bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"),
						event.Click(a.sendMessage),
					),
					vecty.Text("Send"),
				),
			),
		),
	)
}

func main() {
	app := NewApp()
	vecty.SetTitle("Chat Application")
	vecty.RenderBody(app)
}
