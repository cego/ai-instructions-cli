package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cego/ai-instructions/internal/detect"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect project stack from composer.json and package.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := detect.DetectStack(".")
		if err != nil {
			return err
		}

		fmt.Println("Detected stack:")
		if stack.PHP != "" {
			fmt.Printf("- PHP: %s\n", stack.PHP)
		}
		if stack.Laravel != "" {
			fmt.Printf("- Laravel: %s\n", stack.Laravel)
		}
		if stack.Nuxt != "" {
			fmt.Printf("- Nuxt: %s\n", stack.Nuxt)
		}
		if stack.Vue != "" {
			fmt.Printf("- Vue: %s\n", stack.Vue)
		}
		if stack.NuxtUI != "" {
			fmt.Printf("- Nuxt UI: %s\n", stack.NuxtUI)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
