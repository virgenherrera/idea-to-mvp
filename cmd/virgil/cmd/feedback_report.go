package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/virgenherrera/virgil/internal/distribution"
	"github.com/virgenherrera/virgil/internal/feedback"
	"github.com/spf13/cobra"
)

func init() {
	feedbackCmd.AddCommand(reportCmd)
	feedbackCmd.AddCommand(exportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a feedback report",
	RunE: func(cmd *cobra.Command, args []string) error {
		feedbackDir := filepath.Join(distribution.VirgilDir(), "feedback")
		entries, err := feedback.ReadAll(feedbackDir)
		if err != nil {
			return fmt.Errorf("failed to read feedback: %w", err)
		}

		fmt.Print(feedback.GenerateReport(entries))
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all feedback as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		feedbackDir := filepath.Join(distribution.VirgilDir(), "feedback")
		entries, err := feedback.ReadAll(feedbackDir)
		if err != nil {
			return fmt.Errorf("failed to read feedback: %w", err)
		}

		data, err := feedback.Export(entries)
		if err != nil {
			return fmt.Errorf("failed to export feedback: %w", err)
		}

		_, err = os.Stdout.Write(data)
		if err != nil {
			return err
		}
		fmt.Println()
		return nil
	},
}
