package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available rule files",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules := []string{}

		err := filepath.Walk("rules", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				return nil
			}

			trimmed := strings.TrimPrefix(path, "rules/")
			trimmed = strings.TrimSuffix(trimmed, ".md")

			rules = append(rules, trimmed)
			return nil
		})

		if err != nil {
			return err
		}

		if len(rules) == 0 {
			fmt.Println("No rules found.")
			return nil
		}

		for _, rule := range rules {
			fmt.Println(rule)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
