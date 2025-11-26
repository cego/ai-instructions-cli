// go
// File: 'cmd/validate.go'
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cego/ai-instructions/internal/detect"
	"github.com/cego/ai-instructions/rules"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate tech stack and ensure generated files are up to date",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Basic embed sanity check
		list, err := rules.List()
		if err != nil {
			return fmt.Errorf("rules.List failed: %w", err)
		}
		if len(list) == 0 {
			return fmt.Errorf("embedded rules are empty")
		}

		// 2) Detect stack
		stack, err := detect.DetectStack(".")
		if err != nil {
			return fmt.Errorf("stack detection failed: %w", err)
		}

		// Resolve general rules
		generalIDs := buildGeneralRulesFromDetection(stack)
		if len(generalIDs) == 0 {
			return fmt.Errorf("no general rules resolved from detection")
		}
		for _, id := range generalIDs {
			if !ruleExists(id) {
				return fmt.Errorf("missing embedded rule: 'rules/%s.md'", id)
			}
		}
		generalContent, err := loadAndMergeRules(generalIDs)
		if err != nil {
			return fmt.Errorf("failed to merge general rules: %w", err)
		}

		// Prepend stack section like generate does
		stackSection := buildStackSection(stack)
		if stackSection != "" {
			var b bytes.Buffer
			b.WriteString(stackSection)
			b.WriteString("\n\n---\n\n")
			b.WriteString(generalContent)
			generalContent = b.String()
		}

		// Compare current files against expected content
		copilotPath := filepath.ToSlash(".github/copilot-instructions.md")
		agentsPath := filepath.ToSlash("AGENTS.md")

		copilotStatus := compareFileStatus(copilotPath, generalContent)
		agentsStatus := compareFileStatus(agentsPath, generalContent)

		// 5) Report detailed status
		var hadError bool
		switch copilotStatus {
		case statusMissing:
			fmt.Printf("Missing: '%s'\n", copilotPath)
			hadError = true
		case statusOutdated:
			fmt.Printf("Outdated: '%s'\n", copilotPath)
			hadError = true
		case statusUpToDate:
			fmt.Printf("Up to date: '%s'\n", copilotPath)
		}

		switch agentsStatus {
		case statusMissing:
			fmt.Printf("Missing: '%s'\n", agentsPath)
			hadError = true
		case statusOutdated:
			fmt.Printf("Outdated: '%s'\n", agentsPath)
			hadError = true
		case statusUpToDate:
			fmt.Printf("Up to date: '%s'\n", agentsPath)
		}

		if hadError {
			return fmt.Errorf("validation failed")
		}

		fmt.Println("Validation passed: tech stack detected and files are up to date.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

type fileStatus int

const (
	statusUpToDate fileStatus = iota
	statusMissing
	statusOutdated
)

// compareFileStatus returns whether a file is missing, outdated, or up to date.
func compareFileStatus(path string, expected string) fileStatus {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return statusMissing
		}
		// Treat unreadable as outdated
		return statusOutdated
	}
	if bytes.Equal(bytes.TrimSpace(data), bytes.TrimSpace([]byte(expected))) {
		return statusUpToDate
	}
	return statusOutdated
}
