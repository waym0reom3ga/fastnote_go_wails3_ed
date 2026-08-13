// FastNote go_wails3 — document model and shared action layer.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	editorName = "FastNote"
	version    = "1.0.0"
	portID     = "go_wails3"
)

var appExtensions = []string{".md", ".markdown", ".txt"}

type NoteError struct{ msg string }

func (e *NoteError) Error() string { return e.msg }

func noteErrf(format string, a ...any) error {
	return &NoteError{msg: fmt.Sprintf(format, a...)}
}

type Document struct {
	Path  string
	Text  string
	Dirty bool
}

func (d *Document) SetText(text string) {
	if text != d.Text {
		d.Text = text
		d.Dirty = true
	}
}

func (d *Document) Open(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return noteErrf("cannot open %s: %v", path, stripErr(err))
	}
	d.Path = path
	d.Text = string(data)
	d.Dirty = false
	return nil
}

func (d *Document) InsertText(text string) {
	d.Text += text
	d.Dirty = true
}

func (d *Document) Save() (string, error) {
	if d.Path == "" {
		return "", noteErrf("no file name: use save-as (FR-6)")
	}
	if err := d.write(d.Path); err != nil {
		return "", err
	}
	return d.Path, nil
}

func (d *Document) SaveAs(path string) (string, error) {
	if err := d.write(path); err != nil {
		return "", err
	}
	d.Path = path
	return path, nil
}

func (d *Document) write(path string) error {
	if err := os.WriteFile(path, []byte(d.Text), 0o644); err != nil {
		return noteErrf("cannot save %s: %v", path, stripErr(err))
	}
	d.Dirty = false
	return nil
}

type AppState struct {
	Doc       *Document
	NotesDir  string
	SavedOnce bool
}

func NewAppState(notesDir string) *AppState {
	if notesDir == "" {
		notesDir, _ = os.UserHomeDir()
	}
	abs, err := filepath.Abs(notesDir)
	if err != nil {
		abs = notesDir
	}
	return &AppState{Doc: &Document{}, NotesDir: abs}
}

func stripErr(err error) string {
	if err == nil {
		return ""
	}
	if os.IsNotExist(err) {
		return "no such file"
	}
	if os.IsPermission(err) {
		return "permission denied"
	}
	return err.Error()
}

func actionOpen(state *AppState, path string) error { return state.Doc.Open(path) }

func actionInsert(state *AppState, text string) { state.Doc.InsertText(text) }

func actionSave(state *AppState) (string, error) {
	path, err := state.Doc.Save()
	if err == nil {
		state.SavedOnce = true
	}
	return path, err
}

func actionSaveAs(state *AppState, path string) (string, error) {
	path, err := state.Doc.SaveAs(path)
	if err == nil {
		state.SavedOnce = true
	}
	return path, err
}

func actionExportHTML(state *AppState, path, theme string) error {
	return writeHTMLExport(state.Doc.Text, path, theme, "")
}

func actionExportPDF(state *AppState, path string) error {
	return writePDFExport(state.Doc.Text, path)
}

func RunCLIActions(state *AppState, openPath, insert string, doSave bool, exportPath string) error {
	var err error
	if openPath != "" {
		if err = actionOpen(state, openPath); err != nil {
			return err
		}
	}
	if insert != "" {
		actionInsert(state, insert)
	}
	if doSave {
		if _, err = actionSave(state); err != nil {
			return err
		}
	}
	if exportPath != "" {
		if strings.HasSuffix(strings.ToLower(exportPath), ".pdf") {
			err = actionExportPDF(state, exportPath)
		} else {
			err = actionExportHTML(state, exportPath, "light")
		}
	}
	return err
}

func listErr(dir string, err error) error {
	return noteErrf("cannot list %s: %v", dir, stripErr(err))
}
