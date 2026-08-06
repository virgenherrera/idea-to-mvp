package distribution

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

func Init(tier string, methodologyFS embed.FS) (extracted []string, err error) {
	virgilDir := VirgilDir()

	if _, statErr := os.Stat(virgilDir); statErr == nil {
		return nil, fmt.Errorf(".virgil/ already exists — run 'virgil init --force' or delete it manually")
	}

	docsDir := filepath.Join(virgilDir, "docs")
	feedbackDir := filepath.Join(virgilDir, "feedback")

	if mkErr := os.MkdirAll(docsDir, 0755); mkErr != nil {
		return nil, fmt.Errorf("creating docs directory: %w", mkErr)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(virgilDir)
		}
	}()

	if mkErr := os.MkdirAll(feedbackDir, 0755); mkErr != nil {
		return nil, fmt.Errorf("creating feedback directory: %w", mkErr)
	}

	extracted, err = ExtractDocs(docsDir, methodologyFS, tier)
	if err != nil {
		return nil, fmt.Errorf("extracting docs: %w", err)
	}

	cfg := DefaultConfig(tier)
	if err = WriteConfig(virgilDir, cfg); err != nil {
		return nil, fmt.Errorf("writing config: %w", err)
	}

	return extracted, nil
}
