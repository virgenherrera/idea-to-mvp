package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "virgil",
	Short: "Virgil — AI-assisted development methodology",
	Long: `Virgil distributes and manages the Virgil methodology for
AI-assisted software development. It provides methodology docs,
orchestration protocols, and feedback collection.

Three responsibilities:
  1. Methodology provider (source of truth, RAG-able docs)
  2. Mechanism provider (SM behavior, artifactStore, state machine)
  3. Orchestrator protocol (SM orchestrates specialist agents)`,
}

func Execute() error {
	return rootCmd.Execute()
}
