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
			vecty.Class(
				"flex", "flex-col", "gap-6",
				"p-8",
				"bg-white", "dark:bg-gray-800",
				"rounded-lg", "shadow-lg",
				"max-w-md", "mx-auto", "mt-20",
				"transition-colors", "duration-200",
				"border", "border-gray-200", "dark:border-gray-700",
			),
			event.Submit(u.onSubmit).PreventDefault(),
		),
		elem.Heading2(
			vecty.Markup(
				vecty.Class(
					"text-2xl",
					"font-bold",
					"text-center",
					"text-gray-900", "dark:text-white",
				),
			),
			vecty.Text("Enter Your Username"),
		),
		elem.Input(
			vecty.Markup(
				vecty.Class(
					"p-3",
					"border", "border-gray-300", "dark:border-gray-600",
					"rounded-lg",
					"bg-white", "dark:bg-gray-700",
					"text-gray-900", "dark:text-white",
					"placeholder-gray-500", "dark:placeholder-gray-400",
					"focus:ring-2", "focus:ring-blue-500", "dark:focus:ring-blue-400",
					"focus:border-transparent",
					"transition-colors", "duration-200",
				),
				event.Input(u.onInput),
				vecty.Property("type", "text"),
				vecty.Property("value", u.input),
				vecty.Property("placeholder", "Your username..."),
				vecty.Property("required", true),
			),
		),
		elem.Button(
			vecty.Markup(
				vecty.Class(
					"px-6", "py-3",
					"bg-blue-500", "dark:bg-blue-600",
					"hover:bg-blue-600", "dark:hover:bg-blue-700",
					"text-white",
					"font-medium",
					"rounded-lg",
					"transition-colors", "duration-200",
					"focus:outline-none", "focus:ring-2", "focus:ring-blue-500", "dark:focus:ring-blue-400",
				),
				vecty.Property("type", "submit"),
			),
			vecty.Text("Join Chat"),
		),
	)
}
