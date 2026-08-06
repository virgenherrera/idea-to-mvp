package cmd

import (
	"fmt"

	"github.com/virgenherrera/virgil/internal/distribution"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the virgil version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("virgil %s\n", distribution.Version)
	},
}
