package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	virgil "github.com/virgenherrera/virgil"
	"github.com/virgenherrera/virgil/internal/awareness"
	"github.com/virgenherrera/virgil/internal/distribution"
)

var (
	tier  string
	force bool
)

func init() {
	initCmd.Flags().StringVar(&tier, "tier", "standard", "docs tier: minimal, standard, or full")
	initCmd.Flags().BoolVar(&force, "force", false, "remove existing .virgil/ and re-initialize")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Virgil methodology in the current project",
	Long: `Bootstrap the Virgil methodology in the current directory.

Creates .virgil/ with methodology docs, feedback directory, and config.
Injects an awareness block into the project's agent config file
(CLAUDE.md, AGENTS.md, or .cursorrules).

Tiers control how many docs are extracted:
  minimal   5 core docs (overview, glossary, README, artifact schemas)
  standard  Core + execution, behavior, roles, and operation docs
  full      All 31 methodology docs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		validTiers := map[string]bool{"minimal": true, "standard": true, "full": true}
		if !validTiers[tier] {
			return fmt.Errorf("invalid tier %q — must be minimal, standard, or full", tier)
		}

		if force {
			os.RemoveAll(distribution.VirgilDir())
		}

		extracted, err := distribution.Init(tier, virgil.MethodologyFS)
		if err != nil {
			return err
		}

		targetFile, err := awareness.InjectAwarenessBlock(distribution.VirgilDir())
		if err != nil {
			return err
		}

		fmt.Printf("Virgil initialized (%s tier)\n", tier)
		fmt.Printf("  %d docs extracted to .virgil/docs/\n", len(extracted))
		fmt.Printf("  Awareness block injected into %s\n", targetFile)
		return nil
	},
}
