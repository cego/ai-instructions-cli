package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/ai-instructions/internal/detect"
	"github.com/cego/ai-instructions/rules"
)

var (
	flagRules []string
	flagOut   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate copilot-instructions.md and AGENTS.md based on detected stack or explicit flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := "." // kept for future use (detection only)

		var (
			generalRuleIDs []string
			agentRuleIDs   []agentFile
			stack          *detect.DetectedStack
			err            error
		)

		if anyRuleFlagsSet() {
			// Manual mode
			generalRuleIDs = buildGeneralRulesFromFlags()
			agentRuleIDs = buildAgentRulesFromFlags()
		} else {
			// Auto mode
			stack, err = detect.DetectStack(projectRoot)
			if err != nil {
				return err
			}
			generalRuleIDs = buildGeneralRulesFromDetection(stack)
			agentRuleIDs = buildAgentRulesFromDetection(stack)
		}

		// Generate copilot-instructions.md (general rules)
		if len(generalRuleIDs) > 0 {
			content, err := loadAndMergeRules(generalRuleIDs)
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

			// Write same content to AGENTS.md (per original behavior)
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

		// Agents content (separate aggregation)
		if len(agentRuleIDs) > 0 {
			agentContent := buildAgentContent(agentRuleIDs)
			if flagOut == "-" {
				fmt.Println("\n=== (Agents Section) ===")
				fmt.Println(agentContent)
			} else {
				// Append or create AGENTS.md with agent details separated
				// (Optional enhancement: integrate directly above; kept simple)
			}
		}

		if len(generalRuleIDs) == 0 && len(agentRuleIDs) == 0 {
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
	ID    string // rule identifier without prefix & extension (e.g. php/8/agent)
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

// General rules (general.md)
func buildGeneralRulesFromDetection(stack *detect.DetectedStack) []string {
	var ids []string
	addRulesFor(&ids, "php", stack.PHP)
	addRulesFor(&ids, "laravel", stack.Laravel)
	addRulesFor(&ids, "nuxt", stack.Nuxt)
	addRulesFor(&ids, "vue", stack.Vue)
	addRulesFor(&ids, "nuxt_ui", stack.NuxtUI)
	return ids
}

func addRulesFor(ids *[]string, name, version string) {
	if version == "" {
		return
	}

	// Base general
	addIfExists(ids, name+"/general")

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

	// major.minor/general
	if major != "" && minor != "" {
		addIfExists(ids, name+"/"+major+"."+minor+"/general")
	}

	// major/general
	if major != "" {
		addIfExists(ids, name+"/"+major+"/general")
	}
}

func addIfExists(ids *[]string, id string) {
	if ruleExists(id) {
		*ids = append(*ids, id)
	}
}

func buildGeneralRulesFromFlags() []string {
	var ids []string
	for _, r := range flagRules {
		r = filepath.ToSlash(strings.TrimSpace(r))
		if r == "" {
			continue
		}

		// Try as directory: r/general
		dirGeneral := r + "/general"
		if ruleExists(dirGeneral) {
			ids = append(ids, dirGeneral)
			continue
		}

		// Try direct general if user typed framework only
		justGeneral := r + "/general"
		if ruleExists(justGeneral) {
			ids = append(ids, justGeneral)
			continue
		}

		// Accept direct identifier if user passed full (e.g. php/8/general)
		if ruleExists(r) {
			ids = append(ids, r)
		}
	}
	return ids
}

// Agent rules (agent.md)
func buildAgentRulesFromDetection(stack *detect.DetectedStack) []agentFile {
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

	norm := normalizeVersion(version)
	parts := strings.Split(norm, ".")
	major := parts[0]
	minor := ""
	if len(parts) > 1 {
		minor = parts[1]
	}

	// major.minor/agent
	if major != "" && minor != "" {
		id := name + "/" + major + "." + minor + "/agent"
		if ruleExists(id) {
			*files = append(*files, agentFile{Label: label, ID: id})
			return
		}
	}

	// major/agent
	if major != "" {
		id := name + "/" + major + "/agent"
		if ruleExists(id) {
			*files = append(*files, agentFile{Label: label, ID: id})
			return
		}
	}

	// base/agent
	id := name + "/agent"
	if ruleExists(id) {
		*files = append(*files, agentFile{Label: label, ID: id})
	}
}

func buildAgentRulesFromFlags() []agentFile {
	var files []agentFile
	for _, r := range flagRules {
		r = filepath.ToSlash(strings.TrimSpace(r))
		if r == "" {
			continue
		}
		id := r + "/agent"
		label := deriveRuleLabel(id)
		if ruleExists(id) {
			files = append(files, agentFile{Label: label, ID: id})
		}
	}
	return files
}

// Merge general rule contents
func loadAndMergeRules(ids []string) (string, error) {
	var b strings.Builder
	for _, id := range ids {
		data, err := rules.Get(id)
		if err != nil {
			if b.Len() > 0 {
				b.WriteString("\n\n---\n\n")
			}
			b.WriteString("<!-- Missing instructions for ")
			b.WriteString(deriveRuleLabel(id))
			b.WriteString(" (expected file: rules/")
			b.WriteString(id)
			b.WriteString(".md) -->")
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n---\n\n")
		}
		b.WriteString(data)
	}
	return b.String(), nil
}

// Agent content aggregation
func buildAgentContent(files []agentFile) string {
	var b strings.Builder
	b.WriteString("# Agents\n\n")
	b.WriteString("<!-- Generated by ai-instructions. Do not edit manually. -->\n\n---\n\n")

	for i, af := range files {
		if i > 0 {
			b.WriteString("\n\n---\n\n")
		}
		b.WriteString("## ")
		b.WriteString(af.Label)
		b.WriteString("\n\n")

		data, err := rules.Get(af.ID)
		if err != nil {
			b.WriteString("<!-- Missing agent instructions for ")
			b.WriteString(af.Label)
			b.WriteString(" (expected file: rules/")
			b.WriteString(af.ID)
			b.WriteString(".md) -->")
			continue
		}
		b.WriteString(data)
	}
	return b.String()
}

// Existence probe via embedded rules
func ruleExists(id string) bool {
	_, err := rules.Get(id)
	return err == nil
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

func deriveRuleLabel(id string) string {
	id = strings.TrimSuffix(id, "/general")
	id = strings.TrimSuffix(id, "/agent")
	parts := strings.Split(id, "/")
	return strings.Join(parts, " ")
}

func writeFileWithDirs(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := ensureDir(dir); err != nil {
			return err
		}
	}
	return writeFile(path, data)
}

// This ensures the directory exists, even if it's empty.'
func ensureDir(dir string) error {
	return mkdirAll(dir, 0o755)
}

// This is a wrapper around writeFile to allow future abstraction.
func writeFile(path string, data []byte) error {
	return write(path, data, 0o644)
}

// this is a wrapper around os.MkdirAll and os.WriteFile to allow future abstraction.
var (
	mkdirAll = func(path string, perm uint32) error {
		return osMkdirAll(path, perm)
	}
	write = func(name string, data []byte, perm uint32) error {
		return osWriteFile(name, data, perm)
	}
)

// This is a wrapper around os.MkdirAll to allow future abstraction.
func osMkdirAll(path string, perm uint32) error { return os.MkdirAll(path, os.FileMode(perm)) }

// This is a wrapper around os.WriteFile to allow future abstraction.
func osWriteFile(name string, data []byte, perm uint32) error {
	return os.WriteFile(name, data, os.FileMode(perm))
}
