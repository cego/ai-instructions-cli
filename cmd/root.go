package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ai-instructions",
	Short: "AI Instructions CLI for stack detection and config generation",
}

// Execute This is our required entrypoint, for Cobra CLI
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
