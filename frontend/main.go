//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"go-chat/frontend/components"
	"go-chat/frontend/internal"
	"go-chat/frontend/store"

	"log"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// App is the main application component
type App struct {
	vecty.Core
}

// Mount implements the vecty.Mounter interface
func (a *App) Mount() {
	log.Printf("App mounted")
	store.Listeners.Add(a, func() {
		log.Printf("App rerendering due to store change")
		vecty.Rerender(a)
	})
}

// Unmount implements the vecty.Unmounter interface
func (a *App) Unmount() {
	log.Printf("App unmounted")
	store.Listeners.Remove(a)
}

// Render implements the vecty.Component interface
func (a *App) Render() vecty.ComponentOrHTML {
	log.Printf("App rendering, username: %q", store.Username)

	var content vecty.ComponentOrHTML
	if store.Username == "" {
		content = &components.UsernameForm{}
	} else {
		content = components.NewChat()
	}

	return elem.Body(
		vecty.Markup(
			vecty.Class("min-h-screen", "bg-gray-100"),
		),
		elem.Heading1(
			vecty.Markup(
				vecty.Class("text-2xl", "font-bold", "mb-4", "text-center"),
			),
			vecty.Text("Chat App Test"),
		),
		content,
	)
}

func main() {
	fmt.Printf("Starting Chat Application (Build: %s)\n", internal.BuildHash)
	app := &App{}
	vecty.SetTitle("Chat Application")
	vecty.RenderBody(app)
}
