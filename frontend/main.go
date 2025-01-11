//go:build wasm
// +build wasm

package main

import (
	"go-chat/frontend/components"
	"go-chat/frontend/internal"
	"go-chat/frontend/store"
	"log"
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[frontend] ")
}

// App is the main application component
type App struct {
	vecty.Core
}

// Mount implements the vecty.Mounter interface
func (a *App) Mount() {
	log.Printf("🚀 App mounted")
	store.Listeners.Add(a, func() {
		log.Printf("🔄 App rerendering due to store change")
		vecty.Rerender(a)
	})
}

// Unmount implements the vecty.Unmounter interface
func (a *App) Unmount() {
	log.Printf("👋 App unmounted")
	store.Listeners.Remove(a)
}

// Render implements the vecty.Component interface
func (a *App) Render() vecty.ComponentOrHTML {
	log.Printf("🎨 App rendering | username: %q | darkMode: %v", store.Username, store.IsDarkMode)

	var content vecty.ComponentOrHTML
	if store.Username == "" {
		content = &components.UsernameForm{}
	} else {
		content = components.NewChat()
	}

	// Add dark mode class to html element via JavaScript
	doc := js.Global().Get("document")
	html := doc.Get("documentElement")
	classList := html.Get("classList")
	darkClass := "dark"
	hasDark := classList.Call("contains", darkClass).Bool()

	if store.IsDarkMode && !hasDark {
		log.Printf("🌙 Enabling dark mode")
		classList.Call("add", darkClass)
	} else if !store.IsDarkMode && hasDark {
		log.Printf("☀️ Disabling dark mode")
		classList.Call("remove", darkClass)
	}

	baseClasses := []string{
		"min-h-screen",
		"bg-gray-100",
		"dark:bg-gray-900",
		"text-gray-900",
		"dark:text-white",
		"transition-colors",
		"duration-200",
	}

	return elem.Body(
		vecty.Markup(
			vecty.Class(baseClasses...),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("flex", "justify-between", "items-center", "p-4"),
			),
			elem.Heading1(
				vecty.Markup(
					vecty.Class("text-2xl", "font-bold", "mb-0", "dark:text-white"),
				),
				vecty.Text("Chat App Test"),
			),
			&components.DarkModeToggle{},
		),
		content,
	)
}

func main() {
	log.Printf("🎬 Starting Chat Application (Build: %s)", internal.BuildHash)
	app := &App{}
	vecty.SetTitle("Chat Application")
	vecty.RenderBody(app)
}
