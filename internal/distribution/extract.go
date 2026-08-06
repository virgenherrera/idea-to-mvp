package distribution

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var minimalFiles = map[string]bool{
	"docs/en/overview.md":                   true,
	"docs/en/glossary.md":                   true,
	"docs/en/README.md":                     true,
	"docs/en/planning/artifacts/README.md":  true,
	"docs/en/planning/artifacts/schemas.md": true,
}

var standardFiles = map[string]bool{
	"docs/en/echo-system.md":       true,
	"docs/en/artifact-system.md":   true,
	"docs/en/agile-adaptations.md": true,
}

var standardPrefixes = []string{
	"docs/en/execution/",
	"docs/en/planning/behavior/",
	"docs/en/planning/roles/",
	"docs/en/operation/",
}

func shouldInclude(tier, path string) bool {
	if tier == "full" {
		return true
	}
	if minimalFiles[path] {
		return true
	}
	if tier == "standard" {
		if standardFiles[path] {
			return true
		}
		for _, prefix := range standardPrefixes {
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}

func ExtractDocs(destDir string, methodologyFS embed.FS, tier string) ([]string, error) {
	var extracted []string

	err := fs.WalkDir(methodologyFS, "docs/en", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !shouldInclude(tier, path) {
			return nil
		}

		relPath := strings.TrimPrefix(path, "docs/en/")
		destPath := filepath.Join(destDir, relPath)

		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(destDir)) {
			return fmt.Errorf("path traversal detected: %s", relPath)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", relPath, err)
		}

		data, err := methodologyFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", destPath, err)
		}

		extracted = append(extracted, relPath)
		return nil
	})

	return extracted, err
}
