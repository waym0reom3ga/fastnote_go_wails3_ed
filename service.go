// FastNote go_wails3 — Wails service bound to frontend.

package main

import (
	"os"
	"path/filepath"
)

// FastNoteService exposes Go methods to the Wails frontend.
type FastNoteService struct {
	State       *AppState
	Browser     *FileBrowser
	BrowserMode string
	ThemeIndex  int
	PreviewText string
	StatusText  string
}

func NewFastNoteService(state *AppState) *FastNoteService {
	svc := &FastNoteService{
		State: state,
	}
	svc.PreviewText = RenderPlain(state.Doc.Text)
	FnEvent(state, "painted")
	return svc
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
