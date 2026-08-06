package distribution

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

func linkSourcePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".virgil", "link-source"), nil
}

func Register(sourcePath string) error {
	abs, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	docsDir := filepath.Join(abs, "docs", "en")
	info, err := os.Stat(docsDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s does not contain a docs/en/ directory — not a virgil source repo", abs)
	}

	linkFile, err := linkSourcePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(linkFile), 0755); err != nil {
		return fmt.Errorf("cannot create ~/.virgil/: %w", err)
	}

	if err := os.WriteFile(linkFile, []byte(abs), 0644); err != nil {
		return fmt.Errorf("cannot write link-source: %w", err)
	}

	fmt.Printf("Registered source: %s\n", abs)
	return nil
}

func Link() error {
	linkFile, err := linkSourcePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(linkFile)
	if err != nil {
		return fmt.Errorf("no source registered — run 'virgil link --register <path>' first")
	}
	source := string(data)

	virgilDir := VirgilDir()
	if _, err := os.Stat(virgilDir); os.IsNotExist(err) {
		return fmt.Errorf(".virgil/ directory not found — run 'virgil init' first")
	}

	docsDir := filepath.Join(virgilDir, "docs")
	if err := os.RemoveAll(docsDir); err != nil {
		return fmt.Errorf("cannot remove existing docs: %w", err)
	}

	linkTarget := filepath.Join(source, "docs", "en")
	if err := os.Symlink(linkTarget, docsDir); err != nil {
		return fmt.Errorf("cannot create symlink: %w", err)
	}

	cfg, err := ReadConfig(virgilDir)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}
	cfg.DocsSource = "linked"
	cfg.LinkedPath = source
	if err := WriteConfig(virgilDir, cfg); err != nil {
		return fmt.Errorf("cannot update config: %w", err)
	}

	fmt.Printf("Linked .virgil/docs/ → %s\n", linkTarget)
	return nil
}

func Unlink(methodologyFS embed.FS, tier string) error {
	virgilDir := VirgilDir()
	if _, err := os.Stat(virgilDir); os.IsNotExist(err) {
		return fmt.Errorf(".virgil/ directory not found — run 'virgil init' first")
	}

	cfg, err := ReadConfig(virgilDir)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}
	if cfg.DocsSource != "linked" {
		return fmt.Errorf("docs are not linked — nothing to unlink")
	}

	docsDir := filepath.Join(virgilDir, "docs")
	if err := os.RemoveAll(docsDir); err != nil {
		return fmt.Errorf("cannot remove linked docs: %w", err)
	}

	extracted, extErr := ExtractDocs(docsDir, methodologyFS, tier)
	if extErr != nil {
		return fmt.Errorf("cannot re-extract embedded docs: %w", extErr)
	}

	cfg.DocsSource = "embedded"
	cfg.LinkedPath = ""
	if err := WriteConfig(virgilDir, cfg); err != nil {
		return fmt.Errorf("cannot update config: %w", err)
	}

	fmt.Printf("Unlinked. Re-extracted %d embedded docs.\n", len(extracted))
	return nil
}
