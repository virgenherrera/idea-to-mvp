package cmd

import (
	virgil "github.com/virgenherrera/virgil"
	"github.com/virgenherrera/virgil/internal/distribution"
	"github.com/spf13/cobra"
)

var registerPath string

func init() {
	linkCmd.Flags().StringVar(&registerPath, "register", "", "register a virgil source repo path (defaults to current directory)")
	linkCmd.Flags().Lookup("register").NoOptDefVal = "."
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(unlinkCmd)
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link .virgil/docs/ to a local virgil source repo",
	Long: `Without flags: creates a symlink from .virgil/docs/ to the registered source repo.
With --register: registers a virgil source repo path for future linking.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("register") {
			path := registerPath
			if path == "" {
				path = "."
			}
			return distribution.Register(path)
		}
		return distribution.Link()
	},
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove symlink and restore embedded docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := distribution.ReadConfig(distribution.VirgilDir())
		if err != nil {
			return err
		}
		return distribution.Unlink(virgil.MethodologyFS, cfg.Tier)
	},
}
