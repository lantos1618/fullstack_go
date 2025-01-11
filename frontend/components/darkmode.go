//go:build wasm
// +build wasm

package components

import (
	"go-chat/frontend/actions"
	"go-chat/frontend/dispatcher"
	"go-chat/frontend/store"
	"log"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

// DarkModeToggle is a component for toggling dark mode
type DarkModeToggle struct {
	vecty.Core
}

// Mount implements the vecty.Mounter interface
func (d *DarkModeToggle) Mount() {
	log.Printf("🎯 DarkModeToggle mounted")
	store.Listeners.Add(d, func() {
		log.Printf("🔄 DarkModeToggle rerendering | darkMode: %v", store.IsDarkMode)
		vecty.Rerender(d)
	})
}

// Unmount implements the vecty.Unmounter interface
func (d *DarkModeToggle) Unmount() {
	log.Printf("👋 DarkModeToggle unmounted")
	store.Listeners.Remove(d)
}

func (d *DarkModeToggle) onToggle(e *vecty.Event) {
	log.Printf("🔘 Dark mode toggle clicked | current state: %v", store.IsDarkMode)
	dispatcher.Dispatch(&actions.ToggleDarkMode{})
}

// Render implements the vecty.Component interface
func (d *DarkModeToggle) Render() vecty.ComponentOrHTML {
	icon := "🌞" // sun
	if store.IsDarkMode {
		icon = "🌙" // moon
	}

	return elem.Button(
		vecty.Markup(
			vecty.Class("p-2", "rounded-full", "hover:bg-gray-200", "dark:hover:bg-gray-700", "transition-colors", "duration-200", "focus:outline-none", "focus:ring-2", "focus:ring-blue-500"),
			event.Click(d.onToggle),
		),
		vecty.Text(icon),
	)
}
