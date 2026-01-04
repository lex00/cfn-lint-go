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

// SAMTransformOptions configures SAM template transformation.
type SAMTransformOptions struct {
	// Region is the AWS region for transformation context.
	Region string `yaml:"region" json:"region"`

	// AccountID is the AWS account ID for transformation context.
	AccountID string `yaml:"account_id" json:"account_id"`

	// StackName is the CloudFormation stack name for transformation context.
	StackName string `yaml:"stack_name" json:"stack_name"`

	// Partition is the AWS partition (aws, aws-cn, aws-us-gov).
	Partition string `yaml:"partition" json:"partition"`
}

// SAMConfig configures SAM template handling.
type SAMConfig struct {
	// AutoTransform enables automatic SAM to CloudFormation transformation.
	// When true (default), SAM templates are automatically transformed before linting.
	AutoTransform bool `yaml:"auto_transform" json:"auto_transform"`

	// TransformOptions configures the SAM transformation.
	TransformOptions *SAMTransformOptions `yaml:"transform_options" json:"transform_options"`
}

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

	// SAM configures SAM template handling.
	SAM *SAMConfig `yaml:"sam" json:"sam"`
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

	// SAM: override takes precedence if set
	result.SAM = mergeSAMConfig(base.SAM, override.SAM)

	return result
}

// mergeSAMConfig merges two SAM configs, with override taking precedence.
func mergeSAMConfig(base, override *SAMConfig) *SAMConfig {
	if override == nil && base == nil {
		return nil
	}

	if override == nil {
		// Copy base
		result := &SAMConfig{
			AutoTransform: base.AutoTransform,
		}
		if base.TransformOptions != nil {
			result.TransformOptions = &SAMTransformOptions{
				Region:    base.TransformOptions.Region,
				AccountID: base.TransformOptions.AccountID,
				StackName: base.TransformOptions.StackName,
				Partition: base.TransformOptions.Partition,
			}
		}
		return result
	}

	if base == nil {
		// Copy override
		result := &SAMConfig{
			AutoTransform: override.AutoTransform,
		}
		if override.TransformOptions != nil {
			result.TransformOptions = &SAMTransformOptions{
				Region:    override.TransformOptions.Region,
				AccountID: override.TransformOptions.AccountID,
				StackName: override.TransformOptions.StackName,
				Partition: override.TransformOptions.Partition,
			}
		}
		return result
	}

	// Both exist, override takes precedence
	result := &SAMConfig{
		AutoTransform: override.AutoTransform,
	}

	// Merge TransformOptions
	if override.TransformOptions != nil {
		result.TransformOptions = &SAMTransformOptions{
			Region:    override.TransformOptions.Region,
			AccountID: override.TransformOptions.AccountID,
			StackName: override.TransformOptions.StackName,
			Partition: override.TransformOptions.Partition,
		}
		// Fill in missing fields from base if present
		if base.TransformOptions != nil {
			if result.TransformOptions.Region == "" {
				result.TransformOptions.Region = base.TransformOptions.Region
			}
			if result.TransformOptions.AccountID == "" {
				result.TransformOptions.AccountID = base.TransformOptions.AccountID
			}
			if result.TransformOptions.StackName == "" {
				result.TransformOptions.StackName = base.TransformOptions.StackName
			}
			if result.TransformOptions.Partition == "" {
				result.TransformOptions.Partition = base.TransformOptions.Partition
			}
		}
	} else if base.TransformOptions != nil {
		result.TransformOptions = &SAMTransformOptions{
			Region:    base.TransformOptions.Region,
			AccountID: base.TransformOptions.AccountID,
			StackName: base.TransformOptions.StackName,
			Partition: base.TransformOptions.Partition,
		}
	}

	return result
}
