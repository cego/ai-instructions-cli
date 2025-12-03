package rules

import (
	"embed"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Embed the entire rules directory (this directory) recursively.
//
//go:embed **/*.md
var embeddedFS embed.FS

// List returns all markdown rule identifiers (relative path without .md).
func List() ([]string, error) {
	var out []string
	err := fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		// Remove leading "./" if present.
		if strings.HasPrefix(path, "./") {
			path = strings.TrimPrefix(path, "./")
		}
		name := strings.TrimSuffix(path, ".md")
		out = append(out, name)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// Get returns the markdown content for a rule (name is relative path without .md).
func Get(name string) (string, error) {
	data, err := embeddedFS.ReadFile(name + ".md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
