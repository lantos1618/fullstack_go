//go:build wasm
// +build wasm

package components

import (
	"go-chat/frontend/actions"
	"go-chat/frontend/dispatcher"

	"log"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

// UsernameForm is a component for setting the username
type UsernameForm struct {
	vecty.Core
	input string
}

func (u *UsernameForm) onInput(e *vecty.Event) {
	u.input = e.Target.Get("value").String()
	vecty.Rerender(u)
}

func (u *UsernameForm) onSubmit(e *vecty.Event) {
	if u.input == "" {
		return
	}

	log.Printf("Setting username to: %s", u.input)
	dispatcher.Dispatch(&actions.SetUsername{
		Username: u.input,
	})
	vecty.Rerender(u)
}

// Render implements the vecty.Component interface
func (u *UsernameForm) Render() vecty.ComponentOrHTML {
	return elem.Form(
		vecty.Markup(
			vecty.Class("flex", "flex-col", "gap-4", "p-8", "bg-white", "rounded-lg", "shadow-md", "max-w-md", "mx-auto", "mt-20"),
			event.Submit(u.onSubmit).PreventDefault(),
		),
		elem.Heading2(
			vecty.Markup(
				vecty.Class("text-xl", "font-bold", "text-center"),
			),
			vecty.Text("Enter Your Username"),
		),
		elem.Input(
			vecty.Markup(
				vecty.Class("p-2", "border", "rounded"),
				event.Input(u.onInput),
				vecty.Property("type", "text"),
				vecty.Property("value", u.input),
				vecty.Property("placeholder", "Your username..."),
				vecty.Property("required", true),
			),
		),
		elem.Button(
			vecty.Markup(
				vecty.Class("bg-blue-500", "text-white", "px-4", "py-2", "rounded", "hover:bg-blue-600"),
				vecty.Property("type", "submit"),
			),
			vecty.Text("Join Chat"),
		),
	)
}
