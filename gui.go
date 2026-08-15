// FastNote go_wails3 — GUI entry point (separated to avoid WebKit init for --version).

package main

import (
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func runGUI(eventFile string) {
	notesDir, _ := os.UserHomeDir()
	state := NewAppState(notesDir)
	state.EventFile = eventFile
	svc := NewFastNoteService(state)

	app := application.New(application.Options{
		Name:        "FastNote",
		Description: "Markdown Editor",
		Services: []application.Service{
			application.NewService(svc),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "FastNote",
		Width:  1080,
		Height: 740,
		URL:    "/",
	})

	app.Run()
}
