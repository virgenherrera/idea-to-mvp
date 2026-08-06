package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/virgenherrera/virgil/internal/distribution"
	"github.com/virgenherrera/virgil/internal/feedback"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(feedbackCmd)
	feedbackCmd.AddCommand(gapCmd)
	feedbackCmd.AddCommand(frictionCmd)
	feedbackCmd.AddCommand(winCmd)
	feedbackCmd.AddCommand(sessionCmd)

	for _, cmd := range []*cobra.Command{gapCmd, frictionCmd, winCmd} {
		cmd.Flags().String("doc", "", "Document name (required)")
		cmd.Flags().String("section", "", "Section name (required)")
		cmd.Flags().String("scenario", "", "Scenario description (required)")
		cmd.Flags().String("impact", "", "Impact level: low|medium|high (required)")
		cmd.Flags().String("description", "", "Description of the feedback (required)")
		cmd.MarkFlagRequired("doc")
		cmd.MarkFlagRequired("section")
		cmd.MarkFlagRequired("scenario")
		cmd.MarkFlagRequired("impact")
		cmd.MarkFlagRequired("description")
	}

	sessionCmd.Flags().String("goal", "", "Session goal (required)")
	sessionCmd.Flags().String("docs-used", "", "Comma-separated list of docs used (required)")
	sessionCmd.Flags().String("docs-missing", "", "Comma-separated list of missing docs")
	sessionCmd.Flags().Int("gaps-hit", 0, "Number of gaps encountered")
	sessionCmd.Flags().Int("friction-points", 0, "Number of friction points encountered")
	sessionCmd.Flags().String("methodology-grade", "", "Grade: A|B|C|D|F (required)")
	sessionCmd.Flags().String("notes", "", "Additional notes")
	sessionCmd.MarkFlagRequired("goal")
	sessionCmd.MarkFlagRequired("docs-used")
	sessionCmd.MarkFlagRequired("methodology-grade")
}

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Collect and report methodology feedback",
	Long:  "Record gaps, friction points, wins, and session reviews for the Virgil methodology.",
}

var gapCmd = &cobra.Command{
	Use:   "gap",
	Short: "Record a methodology gap",
	RunE:  runFeedbackEntry("gap"),
}

var frictionCmd = &cobra.Command{
	Use:   "friction",
	Short: "Record a friction point",
	RunE:  runFeedbackEntry("friction"),
}

var winCmd = &cobra.Command{
	Use:   "win",
	Short: "Record a methodology win",
	RunE:  runFeedbackEntry("win"),
}

func runFeedbackEntry(entryType string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, feedbackDir, err := loadFeedbackContext()
		if err != nil {
			return err
		}

		impact, _ := cmd.Flags().GetString("impact")
		if impact != "low" && impact != "medium" && impact != "high" {
			return fmt.Errorf("invalid impact %q: must be low, medium, or high", impact)
		}

		doc, _ := cmd.Flags().GetString("doc")
		section, _ := cmd.Flags().GetString("section")
		scenario, _ := cmd.Flags().GetString("scenario")
		description, _ := cmd.Flags().GetString("description")

		entry := feedback.Entry{
			Type:          entryType,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			Doc:           doc,
			Section:       section,
			Scenario:      scenario,
			Impact:        impact,
			Description:   description,
			ProjectTier:   cfg.Tier,
			VirgilVersion: cfg.Version,
		}

		if err := feedback.Write(feedbackDir, entry); err != nil {
			return fmt.Errorf("failed to write feedback: %w", err)
		}

		fmt.Printf("Feedback recorded: %s — %s/%s\n", entryType, doc, section)
		return nil
	}
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Record a session review",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, feedbackDir, err := loadFeedbackContext()
		if err != nil {
			return err
		}

		grade, _ := cmd.Flags().GetString("methodology-grade")
		validGrades := map[string]bool{"A": true, "B": true, "C": true, "D": true, "F": true}
		if !validGrades[grade] {
			return fmt.Errorf("invalid methodology-grade %q: must be A, B, C, D, or F", grade)
		}

		goal, _ := cmd.Flags().GetString("goal")
		docsUsed, _ := cmd.Flags().GetString("docs-used")
		docsMissing, _ := cmd.Flags().GetString("docs-missing")
		gapsHit, _ := cmd.Flags().GetInt("gaps-hit")
		frictionPoints, _ := cmd.Flags().GetInt("friction-points")
		notes, _ := cmd.Flags().GetString("notes")

		entry := feedback.Entry{
			Type:             "session",
			Timestamp:        time.Now().UTC().Format(time.RFC3339),
			Goal:             goal,
			DocsUsed:         docsUsed,
			DocsMissing:      docsMissing,
			GapsHit:          gapsHit,
			FrictionPoints:   frictionPoints,
			MethodologyGrade: grade,
			Notes:            notes,
			ProjectTier:      cfg.Tier,
			VirgilVersion:    cfg.Version,
		}

		if err := feedback.Write(feedbackDir, entry); err != nil {
			return fmt.Errorf("failed to write feedback: %w", err)
		}

		fmt.Printf("Feedback recorded: session — %s\n", goal)
		return nil
	},
}

func loadFeedbackContext() (distribution.Config, string, error) {
	virgilDir := distribution.VirgilDir()
	if _, err := os.Stat(virgilDir); os.IsNotExist(err) {
		return distribution.Config{}, "", fmt.Errorf(".virgil/ directory not found — run 'virgil init' first")
	}

	cfg, err := distribution.ReadConfig(virgilDir)
	if err != nil {
		return distribution.Config{}, "", fmt.Errorf("failed to read config: %w", err)
	}

	feedbackDir := filepath.Join(virgilDir, "feedback")
	return cfg, feedbackDir, nil
}
