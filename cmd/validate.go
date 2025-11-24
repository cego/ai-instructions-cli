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
	validateFile string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate that the generated instructions file is up to date",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := "."

		// 1) Detect stack
		stack, err := detect.DetectStack(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to detect stack: %w", err)
		}

		// 2) Build rule files from detection (auto-mode only)
		files := buildGeneralFilesFromDetection(stack)
		if len(files) == 0 {
			return fmt.Errorf("no rule files selected from detection – nothing to validate")
		}

		// 3) Generate expected content (same as generate auto-mode)
		expected, err := loadAndMergeFiles(files)
		if err != nil {
			return fmt.Errorf("failed to build expected instructions: %w", err)
		}

		stackSection := buildStackSection(stack)
		if stackSection != "" {
			expected = stackSection + "\n\n---\n\n" + expected
		}

		// 4) Load existing file
		data, err := os.ReadFile(validateFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("instructions file %s is missing; run 'ai-instructions generate' first", validateFile)
			}
			return fmt.Errorf("failed to read %s: %w", validateFile, err)
		}

		existing := string(data)

		// 5) Normalize and compare
		if normalizeContent(expected) != normalizeContent(existing) {
			return fmt.Errorf("instructions file %s is outdated; run 'ai-instructions generate' and commit the changes", validateFile)
		}

		// All good
		fmt.Println("✔ Instructions file is up to date.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(
		&validateFile,
		"file",
		filepath.Join(".github", "copilot-instructions.md"),
		"Instructions file to validate",
	)
}

// normalizeContent makes comparison robust to trailing whitespace and line ending differences.
func normalizeContent(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return s
}
