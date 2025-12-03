package cmd

import (
	"fmt"

	"github.com/cego/ai-instructions/rules"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available embedded rule files",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := rules.List()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("No rules found.")
			return nil
		}
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
