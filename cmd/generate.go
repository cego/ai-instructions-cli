package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/ai-instructions/internal/detect"
)

var (
	flagRules []string
	flagOut   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate copilot-instructions.md and AGENTS.md based on detected stack or explicit flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := "."

		var (
			generalFiles []string
			agentFiles   []agentFile
			stack        *detect.DetectedStack
			err          error
		)

		if anyRuleFlagsSet() {
			// Manual mode
			generalFiles = buildGeneralFilesFromFlags()
			agentFiles = buildAgentFilesFromFlags()
		} else {
			// Auto mode
			stack, err = detect.DetectStack(projectRoot)
			if err != nil {
				return err
			}

			generalFiles = buildGeneralFilesFromDetection(stack)
			agentFiles = buildAgentFilesFromDetection(stack)
		}

		// Generate copilot-instructions.md
		if len(generalFiles) > 0 {
			content, err := loadAndMergeFiles(generalFiles)
			if err != nil {
				return err
			}

			// Prepend stack section in auto-mode
			if !anyRuleFlagsSet() {
				stackSection := buildStackSection(stack)
				if stackSection != "" {
					content = stackSection + "\n\n---\n\n" + content
				}
			}

			outPath := flagOut
			if outPath == "" {
				outPath = ".github/copilot-instructions.md"
			}

			if outPath == "-" {
				fmt.Println("=== copilot-instructions.md ===")
				fmt.Println(content)
			} else {
				if err := writeFileWithDirs(outPath, []byte(content)); err != nil {
					return err
				}
				fmt.Printf("Generated instructions\nCOPILOT documentation written to %s\n", outPath)
			}
		}

		// This is where we generate the AGENTS.md
		if len(agentFiles) > 0 {
			content := buildAgentContent(agentFiles)

			// Prepend stack section in auto-mode
			if !anyRuleFlagsSet() && stack != nil {
				stackSection := buildStackSection(stack)
				if stackSection != "" {
					content = stackSection + "\n\n---\n\n" + content
				}
			}

			agentsPath := "AGENTS.md"
			if flagOut == "-" {
				fmt.Println("\n=== AGENTS.md ===")
				fmt.Println(content)
			} else {
				if err := writeFileWithDirs(agentsPath, []byte(content)); err != nil {
					return err
				}
				fmt.Printf("AGENTS documentation written to %s\n", agentsPath)
			}
		}

		if len(generalFiles) == 0 && len(agentFiles) == 0 {
			fmt.Println("No rule files selected – nothing to generate.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringSliceVar(
		&flagRules,
		"rule",
		nil,
		"Rule set(s) to include, e.g. 'php', 'php/8', 'laravel/9'",
	)

	generateCmd.Flags().StringVarP(
		&flagOut,
		"out",
		"o",
		"",
		"Output path for copilot-instructions.md (default .github/copilot-instructions.md, use '-' for stdout)",
	)
}

type agentFile struct {
	Label string
	Path  string
}

func buildStackSection(stack *detect.DetectedStack) string {
	if stack == nil {
		return ""
	}

	var lines []string

	if stack.PHP != "" {
		lines = append(lines, fmt.Sprintf("- PHP: %s", stack.PHP))
	}
	if stack.Laravel != "" {
		lines = append(lines, fmt.Sprintf("- Laravel: %s", stack.Laravel))
	}
	if stack.Nuxt != "" {
		lines = append(lines, fmt.Sprintf("- Nuxt: %s", stack.Nuxt))
	}
	if stack.Vue != "" {
		lines = append(lines, fmt.Sprintf("- Vue: %s", stack.Vue))
	}
	if stack.NuxtUI != "" {
		lines = append(lines, fmt.Sprintf("- Nuxt UI: %s", stack.NuxtUI))
	}

	if len(lines) == 0 {
		return ""
	}

	return "## Stack\n\n" + strings.Join(lines, "\n")
}

func anyRuleFlagsSet() bool {
	return len(flagRules) > 0
}

// Build general.md files
func buildGeneralFilesFromDetection(stack *detect.DetectedStack) []string {
	var files []string

	addRulesFor(&files, "php", stack.PHP)
	addRulesFor(&files, "laravel", stack.Laravel)
	addRulesFor(&files, "nuxt", stack.Nuxt)
	addRulesFor(&files, "vue", stack.Vue)
	addRulesFor(&files, "nuxt_ui", stack.NuxtUI)

	return files
}

func addRulesFor(files *[]string, name, version string) {
	if version == "" {
		return
	}

	baseDir := filepath.Join("rules", name)

	*files = append(*files, filepath.Join(baseDir, "general.md"))

	norm := normalizeVersion(version)
	if norm == "" {
		return
	}

	parts := strings.Split(norm, ".")
	major := parts[0]
	minor := ""
	if len(parts) > 1 {
		minor = parts[1]
	}

	if major != "" && minor != "" {
		dir := filepath.Join(baseDir, major+"."+minor)
		addAllMarkdownInDir(files, dir)
	}

	if major != "" {
		dir := filepath.Join(baseDir, major)
		addAllMarkdownInDir(files, dir)
	}
}

func addAllMarkdownInDir(files *[]string, dir string) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return
	}

	matches, err := filepath.Glob(filepath.Join(dir, "general.md"))
	if err != nil {
		return
	}

	*files = append(*files, matches...)
}

func buildGeneralFilesFromFlags() []string {
	var files []string

	for _, r := range flagRules {
		path := filepath.Join("rules", filepath.FromSlash(r))

		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			matches, err := filepath.Glob(filepath.Join(path, "general.md"))
			if err != nil {
				continue
			}
			files = append(files, matches...)
		} else {
			files = append(files, path)
		}
	}

	return files
}

// Build agent.md files
func buildAgentFilesFromDetection(stack *detect.DetectedStack) []agentFile {
	var files []agentFile

	addAgentFor(&files, "PHP", "php", stack.PHP)
	addAgentFor(&files, "Laravel", "laravel", stack.Laravel)
	addAgentFor(&files, "Nuxt", "nuxt", stack.Nuxt)
	addAgentFor(&files, "Vue", "vue", stack.Vue)
	addAgentFor(&files, "Nuxt UI", "nuxt_ui", stack.NuxtUI)

	return files
}

func addAgentFor(files *[]agentFile, label, name, version string) {
	if version == "" {
		return
	}

	baseDir := filepath.Join("rules", name)
	norm := normalizeVersion(version)
	parts := strings.Split(norm, ".")
	major := parts[0]
	minor := ""
	if len(parts) > 1 {
		minor = parts[1]
	}

	// Try major.minor/agent.md first
	if major != "" && minor != "" {
		path := filepath.Join(baseDir, major+"."+minor, "agent.md")
		*files = append(*files, agentFile{Label: label, Path: path})
		return
	}

	// Fallback: major/agent.md
	if major != "" {
		path := filepath.Join(baseDir, major, "agent.md")
		*files = append(*files, agentFile{Label: label, Path: path})
		return
	}

	// Default: base agent.md
	path := filepath.Join(baseDir, "agent.md")
	*files = append(*files, agentFile{Label: label, Path: path})
}

func buildAgentFilesFromFlags() []agentFile {
	var files []agentFile

	for _, r := range flagRules {
		path := filepath.Join("rules", filepath.FromSlash(r), "agent.md")
		label := deriveRuleNameFromPath(path)
		files = append(files, agentFile{Label: label, Path: path})
	}

	return files
}

func buildAgentContent(files []agentFile) string {
	var b strings.Builder

	b.WriteString("# Agents\n\n")
	b.WriteString("<!-- This file is generated by ai-instructions. Do not edit manually. -->\n\n")
	b.WriteString("---\n\n")

	for i, af := range files {
		if i > 0 {
			b.WriteString("\n\n---\n\n")
		}

		b.WriteString("## ")
		b.WriteString(af.Label)
		b.WriteString("\n\n")

		data, err := os.ReadFile(af.Path)
		if err != nil {
			b.WriteString("<!-- Missing agent instructions for ")
			b.WriteString(af.Label)
			b.WriteString(" (expected file: ")
			b.WriteString(af.Path)
			b.WriteString(") -->")
			continue
		}

		b.Write(data)
	}

	return b.String()
}

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	v = strings.Split(v, "||")[0]
	v = strings.Split(v, " ")[0]
	v = strings.TrimLeft(v, "^~><= ")

	var b strings.Builder
	for _, r := range v {
		if (r >= '0' && r <= '9') || r == '.' {
			b.WriteRune(r)
		} else {
			break
		}
	}

	return b.String()
}

func loadAndMergeFiles(files []string) (string, error) {
	var b strings.Builder

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			if b.Len() > 0 {
				b.WriteString("\n\n---\n\n")
			}

			ruleName := deriveRuleNameFromPath(f)

			b.WriteString("<!-- Missing instructions for ")
			b.WriteString(ruleName)
			b.WriteString(" (expected file: ")
			b.WriteString(f)
			b.WriteString(") -->")

			continue
		}

		if b.Len() > 0 {
			b.WriteString("\n\n---\n\n")
		}

		b.Write(data)
	}

	return b.String(), nil
}

func writeFileWithDirs(path string, data []byte) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(path, data, 0o644)
}

func deriveRuleNameFromPath(path string) string {
	p := filepath.Clean(path)

	prefix := "rules" + string(filepath.Separator)
	p = strings.TrimPrefix(p, prefix)
	p = strings.TrimSuffix(p, ".md")

	parts := strings.Split(p, string(filepath.Separator))
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts, " ")
}
