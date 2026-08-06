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

	globalDocs := GlobalDocsDir()
	if err := os.RemoveAll(globalDocs); err != nil {
		return fmt.Errorf("cannot remove existing docs: %w", err)
	}

	linkTarget := filepath.Join(source, "docs", "en")
	if err := os.Symlink(linkTarget, globalDocs); err != nil {
		return fmt.Errorf("cannot create symlink: %w", err)
	}

	fmt.Printf("Linked %s → %s\n", globalDocs, linkTarget)
	return nil
}

func Unlink(methodologyFS embed.FS, tier string) error {
	globalDocs := GlobalDocsDir()

	info, err := os.Lstat(globalDocs)
	if err != nil {
		return fmt.Errorf("~/.virgil/docs/ not found — nothing to unlink")
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("~/.virgil/docs/ is not a symlink — nothing to unlink")
	}

	if err := os.Remove(globalDocs); err != nil {
		return fmt.Errorf("cannot remove symlink: %w", err)
	}

	extracted, extErr := ExtractDocs(globalDocs, methodologyFS, tier)
	if extErr != nil {
		return fmt.Errorf("cannot re-extract embedded docs: %w", extErr)
	}

	fmt.Printf("Unlinked. Re-extracted %d embedded docs to %s\n", len(extracted), globalDocs)
	return nil
}
