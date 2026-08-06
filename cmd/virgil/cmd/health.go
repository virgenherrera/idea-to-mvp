package cmd

import (
	"github.com/virgenherrera/virgil/internal/health"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(healthCmd)
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show project health metrics (not implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		health.Run()
	},
}
