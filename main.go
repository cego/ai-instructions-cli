package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	templatesDir = "copilot-templates"
	outputPath   = ".github/copilot-instructions.md"
)

var flagToFile = map[string]string{
	"php":     filepath.Join(templatesDir, "php.json"),
	"laravel": filepath.Join(templatesDir, "laravel.json"),
	"react":   filepath.Join(templatesDir, "react.json"),
	"go":      filepath.Join(templatesDir, "go.json"),
}

type Section struct {
	Heading string   `json:"heading"`
	Text    string   `json:"text,omitempty"`
	Bullets []string `json:"bullets,omitempty"`
}

type Template struct {
	Title    string    `json:"title,omitempty"`
	Sections []Section `json:"sections"`
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelperText()
		os.Exit(2)
	}

	var (
		checkMode bool
		templates []Template
		sets      []string
	)

	for _, a := range args {
		switch a {
		case "--validate":
			checkMode = true

		default:
			path, ok := flagToFile[a]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n\n", a)
				printHelperText()
				os.Exit(2)
			}

			info, err := os.Stat(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Missing file for %s: %s\n", a, path)
				os.Exit(2)
			}
			if info.Size() == 0 {
				fmt.Fprintf(os.Stderr, "Template file is empty: %s\n", path)
				os.Exit(2)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read %s: %v\n", path, err)
				os.Exit(1)
			}

			var tpl Template
			if err := json.Unmarshal(data, &tpl); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid JSON in %s: %v\n", path, err)
				os.Exit(1)
			}

			templates = append(templates, tpl)
			base := strings.TrimSuffix(filepath.Base(path), ".json")
			sets = append(sets, base)
		}
	}

	// De-dublicates the flags
	seen := make(map[string]bool, len(sets))
	dedupSets := make([]string, 0, len(sets))
	dedupTemplates := make([]Template, 0, len(templates))
	for i, s := range sets {
		if seen[s] {
			continue
		}
		seen[s] = true
		dedupSets = append(dedupSets, s)
		dedupTemplates = append(dedupTemplates, templates[i])
	}
	sets = dedupSets
	templates = dedupTemplates

	if len(sets) == 0 {
		printHelperText()
		os.Exit(2)
	}

	title, merged := mergeTemplates(templates, sets)
	md := renderMarkdown(title, merged, sets)

	if checkMode {
		existing, err := os.ReadFile(outputPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			os.Exit(1)
		}
		if string(existing) != md {
			fmt.Fprintf(os.Stderr, "%s is out of date. Re-run without --validate to regenerate.\n", outputPath)
			os.Exit(1)
		}
		fmt.Println("OK: file is up to date.")
		return
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputPath, []byte(md), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %s (sets: %s)\n", outputPath, strings.Join(sets, ", "))
}

func printHelperText() {
	fmt.Println("Usage:")
	fmt.Println("  ai-instructions-pilot [flags]\n")
	fmt.Println("Description:")
	fmt.Println("  Generates .github/copilot-instructions.md from selected templates.\n")
	fmt.Println("Flags:")

	// Sort keys for deterministic help output
	keys := make([]string, 0, len(flagToFile))
	for k := range flagToFile {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, flag := range keys {
		fmt.Printf("  %-10s %s\n", flag, flagToFile[flag])
	}
	fmt.Println("  --validate     validation mode: verify file is up to date, don't write\n")

	fmt.Println("Examples:")
	fmt.Println("  ai-instructions-pilot --laravel")
	fmt.Println("  ai-instructions-pilot --php --react")
	fmt.Println("  ai-instructions-pilot --php --react --validate")
}

func loadTemplate(set string) (Template, error) {
	path := filepath.Join(templatesDir, set+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Template{}, fmt.Errorf("cannot read %s: %w", path, err)
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return Template{}, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return t, nil
}

func mergeTemplates(templates []Template, sets []string) (string, []Section) {
	// Title: if any template has a title, use the first; otherwise compose.
	title := ""
	for _, t := range templates {
		if strings.TrimSpace(t.Title) != "" {
			title = t.Title
			break
		}
	}
	if title == "" {
		title = "Copilot instructions (" + strings.Join(sets, " + ") + ")"
	}

	// Merge sections by heading, preserving first-seen order.
	var merged []Section
	index := make(map[string]int) // heading -> idx
	for _, t := range templates {
		for _, s := range t.Sections {
			h := strings.TrimSpace(s.Heading)
			if h == "" {
				continue
			}
			if idx, ok := index[h]; ok {
				// Merge into existing.
				m := merged[idx]
				// Prefer the first non-empty text; if both non-empty and different, append.
				if strings.TrimSpace(s.Text) != "" {
					if strings.TrimSpace(m.Text) == "" {
						m.Text = s.Text
					} else if strings.TrimSpace(m.Text) != strings.TrimSpace(s.Text) {
						m.Text = strings.TrimSpace(m.Text) + "\n\n" + strings.TrimSpace(s.Text)
					}
				}
				// Merge bullets with de-dupe (stable).
				m.Bullets = mergeBullets(m.Bullets, s.Bullets)
				merged[idx] = m
			} else {
				// Normalize bullets (trim) and add.
				s.Bullets = trimAll(s.Bullets)
				index[h] = len(merged)
				merged = append(merged, s)
			}
		}
	}

	return title, merged
}

func trimAll(in []string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func mergeBullets(a, b []string) []string {
	// stable de-dupe as we append b to a
	seen := make(map[string]bool, len(a)+len(b))
	out := make([]string, 0, len(a)+len(b))
	for _, v := range append(a, b...) {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

func renderMarkdown(title string, sections []Section, sets []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", strings.TrimSpace(title))
	// Keep this comment deterministic to preserve idempotency.
	fmt.Fprintf(&b, "<!-- generated by ai-instructions-pilot: sets: %s -->\n\n", strings.Join(sets, ", "))

	for i, s := range sections {
		fmt.Fprintf(&b, "## %s\n\n", s.Heading)
		if strings.TrimSpace(s.Text) != "" {
			fmt.Fprintf(&b, "%s\n\n", strings.TrimSpace(s.Text))
		}
		for _, bullet := range s.Bullets {
			fmt.Fprintf(&b, "- %s\n", bullet)
		}
		if len(s.Bullets) > 0 {
			b.WriteString("\n")
		}
		// Avoid trailing extra newline at EOF; harmless either way.
		if i == len(sections)-1 {
			// no-op
		}
	}
	return b.String()
}
