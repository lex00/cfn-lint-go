// Package config provides configuration file support for cfn-lint.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the cfn-lint configuration.
type Config struct {
	// Templates is a list of template files or glob patterns to lint.
	Templates []string `yaml:"templates" json:"templates"`

	// IgnoreTemplates is a list of templates or patterns to ignore.
	IgnoreTemplates []string `yaml:"ignore_templates" json:"ignore_templates"`

	// Regions is a list of AWS regions to validate against.
	Regions []string `yaml:"regions" json:"regions"`

	// IgnoreChecks is a list of rule IDs to ignore.
	IgnoreChecks []string `yaml:"ignore_checks" json:"ignore_checks"`

	// IncludeChecks is a list of rule IDs to include (even if ignored).
	IncludeChecks []string `yaml:"include_checks" json:"include_checks"`

	// IncludeExperimental enables experimental rules.
	IncludeExperimental bool `yaml:"include_experimental" json:"include_experimental"`

	// ConfigureRules contains rule-specific configuration.
	ConfigureRules map[string]map[string]interface{} `yaml:"configure_rules" json:"configure_rules"`

	// Format is the output format (text, json, sarif, junit, pretty).
	Format string `yaml:"format" json:"format"`

	// OutputFile is the file to write output to.
	OutputFile string `yaml:"output_file" json:"output_file"`
}

// ConfigFileNames lists the config file names to search for, in order of preference.
var ConfigFileNames = []string{
	".cfnlintrc",
	".cfnlintrc.yaml",
	".cfnlintrc.yml",
	".cfnlintrc.json",
}

// Load loads a configuration file from the specified path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	ext := filepath.Ext(path)
	if ext == "" {
		// Try to detect format from the file
		ext = detectFormat(data)
	}

	var cfg Config

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing YAML config: %w", err)
		}
	default:
		// Try YAML first, then JSON
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			if err := json.Unmarshal(data, &cfg); err != nil {
				return nil, fmt.Errorf("parsing config (tried both YAML and JSON): %w", err)
			}
		}
	}

	return &cfg, nil
}

// Find searches for a config file starting from the current directory
// and walking up to the git root or filesystem root.
func Find() (string, error) {
	// Start from current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Find git root
	gitRoot := findGitRoot(currentDir)

	// Search from current directory up to git root (or filesystem root)
	dir := currentDir
	for {
		// Try each config file name
		for _, name := range ConfigFileNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		// If we've reached git root or filesystem root, stop
		parent := filepath.Dir(dir)
		if parent == dir || (gitRoot != "" && dir == gitRoot) {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no config file found")
}

// findGitRoot finds the git repository root by looking for .git directory.
func findGitRoot(startDir string) string {
	dir := startDir
	for {
		gitDir := filepath.Join(dir, ".git")
		if stat, err := os.Stat(gitDir); err == nil && stat.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// detectFormat tries to detect if data is JSON or YAML.
func detectFormat(data []byte) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(string(data))

	// JSON starts with { or [
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return ".json"
	}

	// Default to YAML
	return ".yaml"
}

// Merge merges two configs, with the override config taking precedence.
func Merge(base, override *Config) *Config {
	result := &Config{}

	// Templates: append both
	result.Templates = append(result.Templates, base.Templates...)
	result.Templates = append(result.Templates, override.Templates...)

	// IgnoreTemplates: append both
	result.IgnoreTemplates = append(result.IgnoreTemplates, base.IgnoreTemplates...)
	result.IgnoreTemplates = append(result.IgnoreTemplates, override.IgnoreTemplates...)

	// Regions: override takes precedence if set
	if len(override.Regions) > 0 {
		result.Regions = override.Regions
	} else {
		result.Regions = base.Regions
	}

	// IgnoreChecks: append both
	result.IgnoreChecks = append(result.IgnoreChecks, base.IgnoreChecks...)
	result.IgnoreChecks = append(result.IgnoreChecks, override.IgnoreChecks...)

	// IncludeChecks: append both
	result.IncludeChecks = append(result.IncludeChecks, base.IncludeChecks...)
	result.IncludeChecks = append(result.IncludeChecks, override.IncludeChecks...)

	// IncludeExperimental: override takes precedence
	result.IncludeExperimental = base.IncludeExperimental || override.IncludeExperimental

	// ConfigureRules: merge maps
	result.ConfigureRules = make(map[string]map[string]interface{})
	for k, v := range base.ConfigureRules {
		result.ConfigureRules[k] = v
	}
	for k, v := range override.ConfigureRules {
		result.ConfigureRules[k] = v
	}

	// Format: override takes precedence if set
	if override.Format != "" {
		result.Format = override.Format
	} else {
		result.Format = base.Format
	}

	// OutputFile: override takes precedence if set
	if override.OutputFile != "" {
		result.OutputFile = override.OutputFile
	} else {
		result.OutputFile = base.OutputFile
	}

	return result
}
