// In-app file browser (spec §3) — toolkit-independent state machine.

package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	Name  string
	IsDir bool
}

type FileBrowser struct {
	Mode      string
	Cwd       string
	ShowAll   bool
	PathInput string
	Selected  string
	Entries   []Entry
}

func NewFileBrowser(mode, startDir string) (*FileBrowser, error) {
	if startDir == "" {
		startDir, _ = os.UserHomeDir()
	}
	abs, err := filepath.Abs(startDir)
	if err != nil {
		return nil, noteErrf("bad start dir %s", startDir)
	}
	b := &FileBrowser{Mode: mode, Cwd: abs}
	if err := b.Refresh(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *FileBrowser) Refresh() error {
	names, err := os.ReadDir(b.Cwd)
	if err != nil {
		return listErr(b.Cwd, err)
	}
	sort.Slice(names, func(i, j int) bool { return names[i].Name() < names[j].Name() })
	entries := []Entry{{Name: "..", IsDir: true}}
	for _, de := range names {
		if de.IsDir() {
			entries = append(entries, Entry{Name: de.Name(), IsDir: true})
		} else if info, err := de.Info(); err == nil && info.Mode().IsRegular() {
			if b.ShowAll || hasAppExt(de.Name()) {
				entries = append(entries, Entry{Name: de.Name()})
			}
		}
	}
	b.Entries = entries
	return nil
}

func hasAppExt(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, e := range appExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func (b *FileBrowser) Activate(name string) (string, error) {
	full := filepath.Join(b.Cwd, name)
	info, err := os.Stat(full)
	if err == nil && info.IsDir() {
		b.Cwd = full
		b.PathInput = ""
		b.Selected = ""
		if err := b.Refresh(); err != nil {
			return "", err
		}
		return "", nil
	}
	if b.Mode == "open" && (err != nil || !info.Mode().IsRegular()) {
		return "", nil
	}
	b.Selected = full
	return full, nil
}

func (b *FileBrowser) Parent() error {
	parent := filepath.Dir(b.Cwd)
	if parent != "" && parent != b.Cwd {
		b.Cwd = parent
		return b.Refresh()
	}
	return nil
}

func (b *FileBrowser) ToggleFilter() error {
	b.ShowAll = !b.ShowAll
	return b.Refresh()
}

func (b *FileBrowser) Result() (string, error) {
	path := strings.TrimSpace(b.PathInput)
	if path == "" {
		path = b.Selected
	}
	if path == "" {
		return "", noteErrf("choose a file or type a path")
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(b.Cwd, path)
	}
	return filepath.Abs(path)
}

func (b *FileBrowser) SelectByPath(path string) {
	b.PathInput = path
	b.Selected = path
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
