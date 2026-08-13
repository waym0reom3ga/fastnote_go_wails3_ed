// FastNote go_wails3 — Wails service bound to frontend.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FastNoteService exposes Go methods to the Wails frontend.
type FastNoteService struct {
	State          *AppState
	Browser        *FileBrowser
	BrowserMode    string
	ThemeIndex     int
	PreviewText    string
	StatusText     string
	controlMapPath string
	readyFilePath  string
}

func NewFastNoteService(state *AppState, openPath, controlMapPath, readyFilePath string) *FastNoteService {
	svc := &FastNoteService{
		State:          state,
		controlMapPath: controlMapPath,
		readyFilePath:  readyFilePath,
	}
	if openPath != "" {
		actionOpen(state, openPath)
	}
	svc.PreviewText = RenderPlain(state.Doc.Text)
	svc.writeControlMap()
	svc.signalReady()
	return svc
}

func (s *FastNoteService) writeControlMap() {
	if s.controlMapPath == "" {
		return
	}
	var buf strings.Builder
	buf.WriteString("name\tx\ty\tw\th\n")
	tb := float32(34)
	controls := []struct {
		name string
		x0, y0, x1, y1 float32
	}{
		{"Open", 6, 6, 74, tb - 6},
		{"Save", 80, 6, 148, tb - 6},
		{"SaveAs", 154, 6, 222, tb - 6},
		{"Export", 228, 6, 296, tb - 6},
		{"ExportPdf", 302, 6, 378, tb - 6},
		{"Theme", 384, 6, 452, tb - 6},
		{"editor", 0, tb, 540, 700},
	}
	for _, c := range controls {
		buf.WriteString(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\n",
			c.name, int(c.x0), int(c.y0), int(c.x1-c.x0), int(c.y1-c.y0)))
	}
	os.WriteFile(s.controlMapPath, []byte(buf.String()), 0o644)
}

func (s *FastNoteService) signalReady() {
	if s.readyFilePath == "" {
		return
	}
	os.WriteFile(s.readyFilePath, nil, 0o644)
}

// GetState returns the current document state to the frontend.
func (s *FastNoteService) GetState() map[string]any {
	return map[string]any{
		"text":       s.State.Doc.Text,
		"path":       s.State.Doc.Path,
		"dirty":      s.State.Doc.Dirty,
		"preview":    s.PreviewText,
		"status":     s.StatusText,
		"themeIndex": s.ThemeIndex,
	}
}

// SetText updates the document content from the frontend editor.
func (s *FastNoteService) SetText(text string) string {
	s.State.Doc.SetText(text)
	s.PreviewText = RenderPlain(s.State.Doc.Text)
	s.StatusText = "Editing"
	return s.PreviewText
}

// OpenFile opens a file by path.
func (s *FastNoteService) OpenFile(path string) map[string]any {
	if err := actionOpen(s.State, path); err != nil {
		return map[string]any{"error": err.Error()}
	}
	s.PreviewText = RenderPlain(s.State.Doc.Text)
	s.StatusText = "Opened " + filepath.Base(path)
	return map[string]any{
		"text":    s.State.Doc.Text,
		"path":    s.State.Doc.Path,
		"preview": s.PreviewText,
		"status":  s.StatusText,
	}
}

// SaveFile saves the current document.
func (s *FastNoteService) SaveFile() map[string]any {
	if s.State.Doc.Path == "" {
		return map[string]any{"error": "no file name: use save-as"}
	}
	path, err := actionSave(s.State)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	s.StatusText = "Saved"
	return map[string]any{"path": path, "status": s.StatusText}
}

// SaveAsFile saves to a new path.
func (s *FastNoteService) SaveAsFile(path string) map[string]any {
	path = ensureNewPath(path)
	p, err := actionSaveAs(s.State, path)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	s.StatusText = "Saved as " + filepath.Base(path)
	return map[string]any{"path": p, "status": s.StatusText}
}

// ExportHTML exports the document as HTML.
func (s *FastNoteService) ExportHTML(path string) map[string]any {
	if s.State.Doc.Path == "" {
		return map[string]any{"error": "open a document first"}
	}
	if err := actionExportHTML(s.State, path, s.currentTheme()); err != nil {
		return map[string]any{"error": err.Error()}
	}
	s.StatusText = "Exported " + filepath.Base(path)
	return map[string]any{"path": path, "status": s.StatusText}
}

// ExportPDF exports the document as PDF.
func (s *FastNoteService) ExportPDF(path string) map[string]any {
	if s.State.Doc.Path == "" {
		return map[string]any{"error": "open a document first"}
	}
	if err := actionExportPDF(s.State, path); err != nil {
		return map[string]any{"error": err.Error()}
	}
	s.StatusText = "Exported " + filepath.Base(path)
	return map[string]any{"path": path, "status": s.StatusText}
}

// ToggleTheme switches between light and dark.
func (s *FastNoteService) ToggleTheme() map[string]any {
	themes := []string{"light", "dark"}
	s.ThemeIndex = (s.ThemeIndex + 1) % len(themes)
	s.StatusText = "Theme: " + themes[s.ThemeIndex]
	return map[string]any{
		"theme":      themes[s.ThemeIndex],
		"themeIndex": s.ThemeIndex,
		"status":     s.StatusText,
	}
}

func (s *FastNoteService) currentTheme() string {
	themes := []string{"light", "dark"}
	return themes[s.ThemeIndex]
}

// ListDir lists directory contents for the in-app browser.
func (s *FastNoteService) ListDir(dir string) map[string]any {
	if dir == "" {
		dir = s.State.NotesDir
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	result := make([]map[string]any, 0, len(entries))
	for _, e := range entries {
		info, _ := e.Info()
		result = append(result, map[string]any{
			"name":  e.Name(),
			"isDir": e.IsDir(),
			"size":  info.Size(),
		})
	}
	return map[string]any{"dir": dir, "entries": result}
}

// GetRenderedFragment returns the rendered HTML fragment for preview.
func (s *FastNoteService) GetRenderedFragment() string {
	baseDir := filepath.Dir(s.State.Doc.Path)
	return RenderFragment(s.State.Doc.Text, baseDir)
}
