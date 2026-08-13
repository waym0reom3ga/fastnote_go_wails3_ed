// FastNote go_wails3 — Wails v3 application entry point.

package main

import (
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	openPath := ""
	controlMapPath := ""
	readyFilePath := ""
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--open":
			if i+1 < len(args) {
				openPath = args[i+1]
				i++
			}
		case "--control-map":
			if i+1 < len(args) {
				controlMapPath = args[i+1]
				i++
			}
		case "--ready-file":
			if i+1 < len(args) {
				readyFilePath = args[i+1]
				i++
			}
		case "--cli":
			runCLI(args[i+1:])
			return
		case "--selftest":
			if RunSelfTest() {
				os.Exit(0)
			}
			os.Exit(1)
		}
	}

	notesDir, _ := os.UserHomeDir()
	state := NewAppState(notesDir)
	svc := NewFastNoteService(state, openPath, controlMapPath, readyFilePath)

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

func runCLI(args []string) {
	notesDir, _ := os.UserHomeDir()
	state := NewAppState(notesDir)
	openPath := ""
	insert := ""
	doSave := false
	exportPath := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--open":
			if i+1 < len(args) {
				openPath = args[i+1]
				i++
			}
		case "--insert":
			if i+1 < len(args) {
				insert = args[i+1]
				i++
			}
		case "--save":
			doSave = true
		case "--export":
			if i+1 < len(args) {
				exportPath = args[i+1]
				i++
			}
		}
	}
	if err := RunCLIActions(state, openPath, insert, doSave, exportPath); err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}
}
