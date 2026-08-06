package distribution

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type FeedbackConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"`
	Path    string `yaml:"path"`
}

type TelemetryConfig struct {
	Enabled bool `yaml:"enabled"`
}

type Config struct {
	Version    string          `yaml:"version"`
	Tier       string          `yaml:"tier"`
	Language   string          `yaml:"language"`
	DocsSource string          `yaml:"docs_source"`
	LinkedPath string          `yaml:"linked_path,omitempty"`
	Feedback   FeedbackConfig  `yaml:"feedback"`
	Telemetry  TelemetryConfig `yaml:"telemetry"`
}

var Version = "0.5.0-dev"

func DefaultConfig(tier string) Config {
	return Config{
		Version:    Version,
		Tier:       tier,
		Language:   "en",
		DocsSource: "embedded",
		Feedback: FeedbackConfig{
			Enabled: true,
			Format:  "jsonl",
			Path:    ".virgil/feedback/",
		},
		Telemetry: TelemetryConfig{
			Enabled: false,
		},
	}
}

func WriteConfig(dir string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0644)
}

func ReadConfig(dir string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	validTiers := map[string]bool{"minimal": true, "standard": true, "full": true}
	if !validTiers[c.Tier] {
		return fmt.Errorf("invalid tier in config: %q — must be minimal, standard, or full", c.Tier)
	}
	validSources := map[string]bool{"embedded": true, "linked": true}
	if !validSources[c.DocsSource] {
		return fmt.Errorf("invalid docs_source in config: %q", c.DocsSource)
	}
	return nil
}

func VirgilDir() string {
	return ".virgil"
}

func ConfigPath() string {
	return filepath.Join(VirgilDir(), "config.yaml")
}
