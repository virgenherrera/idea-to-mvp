package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Assess methodology compliance (not implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("virgil scan — not implemented yet")
		fmt.Println()
		fmt.Println("This command will scan the project for methodology compliance:")
		fmt.Println("  - Artifact presence and completeness")
		fmt.Println("  - Binding coverage (code ↔ artifacts)")
		fmt.Println("  - Gate compliance")
		fmt.Println()
		fmt.Println("See .virgil/docs/overview.md for the full methodology.")
	},
}
