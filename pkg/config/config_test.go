package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_YAML(t *testing.T) {
	content := `
templates:
  - "*.yaml"
  - "templates/*.yaml"
ignore_templates:
  - "test/*.yaml"
regions:
  - us-east-1
  - us-west-2
ignore_checks:
  - W2001
  - W3001
include_checks:
  - E9999
include_experimental: true
format: json
output_file: results.json
configure_rules:
  E3012:
    strict: true
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(cfg.Templates))
	}

	if len(cfg.IgnoreTemplates) != 1 {
		t.Errorf("Expected 1 ignore template, got %d", len(cfg.IgnoreTemplates))
	}

	if len(cfg.Regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(cfg.Regions))
	}

	if len(cfg.IgnoreChecks) != 2 {
		t.Errorf("Expected 2 ignore checks, got %d", len(cfg.IgnoreChecks))
	}

	if len(cfg.IncludeChecks) != 1 {
		t.Errorf("Expected 1 include check, got %d", len(cfg.IncludeChecks))
	}

	if !cfg.IncludeExperimental {
		t.Error("Expected IncludeExperimental to be true")
	}

	if cfg.Format != "json" {
		t.Errorf("Expected format 'json', got %s", cfg.Format)
	}

	if cfg.OutputFile != "results.json" {
		t.Errorf("Expected output_file 'results.json', got %s", cfg.OutputFile)
	}

	if cfg.ConfigureRules["E3012"] == nil {
		t.Error("Expected E3012 rule config to exist")
	}
}

func TestLoad_JSON(t *testing.T) {
	content := `{
  "templates": ["*.json"],
  "regions": ["eu-west-1"],
  "ignore_checks": ["E1001"],
  "format": "sarif"
}`
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(cfg.Templates))
	}

	if cfg.Format != "sarif" {
		t.Errorf("Expected format 'sarif', got %s", cfg.Format)
	}
}

func TestLoad_NonExistent(t *testing.T) {
	_, err := Load("/non/existent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	content := `
templates: [
  - invalid yaml
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	content := `{"templates": invalid}`
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestLoad_NoExtension(t *testing.T) {
	content := `templates:
  - "*.yaml"
`
	tmpFile, err := os.CreateTemp("", "cfnlintrc")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(cfg.Templates))
	}
}

func TestMerge(t *testing.T) {
	base := &Config{
		Templates:           []string{"base/*.yaml"},
		IgnoreTemplates:     []string{"base/ignore/*.yaml"},
		Regions:             []string{"us-east-1"},
		IgnoreChecks:        []string{"W1001"},
		IncludeChecks:       []string{"E1001"},
		IncludeExperimental: false,
		ConfigureRules: map[string]map[string]interface{}{
			"E3012": {"strict": false},
		},
		Format:     "text",
		OutputFile: "base-output.txt",
	}

	override := &Config{
		Templates:           []string{"override/*.yaml"},
		IgnoreTemplates:     []string{"override/ignore/*.yaml"},
		Regions:             []string{"eu-west-1"},
		IgnoreChecks:        []string{"W2001"},
		IncludeChecks:       []string{"E2001"},
		IncludeExperimental: true,
		ConfigureRules: map[string]map[string]interface{}{
			"E3012": {"strict": true},
			"W1001": {"enabled": false},
		},
		Format:     "json",
		OutputFile: "override-output.json",
	}

	result := Merge(base, override)

	// Templates should be combined
	if len(result.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(result.Templates))
	}

	// IgnoreTemplates should be combined
	if len(result.IgnoreTemplates) != 2 {
		t.Errorf("Expected 2 ignore templates, got %d", len(result.IgnoreTemplates))
	}

	// Regions should be from override
	if len(result.Regions) != 1 || result.Regions[0] != "eu-west-1" {
		t.Errorf("Expected regions from override, got %v", result.Regions)
	}

	// IgnoreChecks should be combined
	if len(result.IgnoreChecks) != 2 {
		t.Errorf("Expected 2 ignore checks, got %d", len(result.IgnoreChecks))
	}

	// IncludeChecks should be combined
	if len(result.IncludeChecks) != 2 {
		t.Errorf("Expected 2 include checks, got %d", len(result.IncludeChecks))
	}

	// IncludeExperimental should be true (override OR base)
	if !result.IncludeExperimental {
		t.Error("Expected IncludeExperimental to be true")
	}

	// ConfigureRules should be merged
	if len(result.ConfigureRules) != 2 {
		t.Errorf("Expected 2 rule configs, got %d", len(result.ConfigureRules))
	}

	// Format should be from override
	if result.Format != "json" {
		t.Errorf("Expected format 'json', got %s", result.Format)
	}

	// OutputFile should be from override
	if result.OutputFile != "override-output.json" {
		t.Errorf("Expected output file from override, got %s", result.OutputFile)
	}
}

func TestMerge_EmptyOverride(t *testing.T) {
	base := &Config{
		Templates: []string{"base/*.yaml"},
		Regions:   []string{"us-east-1"},
		Format:    "text",
	}

	override := &Config{}

	result := Merge(base, override)

	if len(result.Regions) != 1 || result.Regions[0] != "us-east-1" {
		t.Errorf("Expected base regions when override is empty")
	}

	if result.Format != "text" {
		t.Errorf("Expected base format when override is empty")
	}
}

func TestFind_NotFound(t *testing.T) {
	// Create a temp directory without config files
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldWd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatalf("Failed to get working directory: %v", wdErr)
	}
	defer os.Chdir(oldWd)

	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("Failed to change directory: %v", chdirErr)
	}

	_, err = Find()
	if err == nil {
		t.Error("Expected error when no config file found")
	}
}

func TestFind_Found(t *testing.T) {
	// Create a temp directory with a config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Resolve symlinks for macOS where /var -> /private/var
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("Failed to eval symlinks: %v", err)
	}

	configPath := filepath.Join(tmpDir, ".cfnlintrc")
	if writeErr := os.WriteFile(configPath, []byte("templates:\n  - \"*.yaml\"\n"), 0644); writeErr != nil {
		t.Fatalf("Failed to write config file: %v", writeErr)
	}

	oldWd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatalf("Failed to get working directory: %v", wdErr)
	}
	defer os.Chdir(oldWd)

	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("Failed to change directory: %v", chdirErr)
	}

	foundPath, findErr := Find()
	if findErr != nil {
		t.Fatalf("Find failed: %v", findErr)
	}

	if foundPath != configPath {
		t.Errorf("Expected path %s, got %s", configPath, foundPath)
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{"JSON object", `{"key": "value"}`, ".json"},
		{"JSON array", `["item1", "item2"]`, ".json"},
		{"JSON with whitespace", `  {"key": "value"}`, ".json"},
		{"YAML", "key: value", ".yaml"},
		{"YAML with dashes", "---\nkey: value", ".yaml"},
		{"Empty", "", ".yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectFormat([]byte(tt.data))
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestConfigFileNames(t *testing.T) {
	if len(ConfigFileNames) == 0 {
		t.Error("ConfigFileNames should not be empty")
	}

	// First should be .cfnlintrc
	if ConfigFileNames[0] != ".cfnlintrc" {
		t.Errorf("First config file name should be .cfnlintrc, got %s", ConfigFileNames[0])
	}
}

func TestLoad_SAMConfig(t *testing.T) {
	content := `
templates:
  - "*.yaml"
sam:
  auto_transform: true
  transform_options:
    region: us-west-2
    account_id: "987654321098"
    stack_name: my-sam-app
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.SAM == nil {
		t.Fatal("Expected SAM config to be non-nil")
	}

	if !cfg.SAM.AutoTransform {
		t.Error("Expected AutoTransform to be true")
	}

	if cfg.SAM.TransformOptions == nil {
		t.Fatal("Expected TransformOptions to be non-nil")
	}

	if cfg.SAM.TransformOptions.Region != "us-west-2" {
		t.Errorf("Expected region 'us-west-2', got %s", cfg.SAM.TransformOptions.Region)
	}

	if cfg.SAM.TransformOptions.AccountID != "987654321098" {
		t.Errorf("Expected account_id '987654321098', got %s", cfg.SAM.TransformOptions.AccountID)
	}

	if cfg.SAM.TransformOptions.StackName != "my-sam-app" {
		t.Errorf("Expected stack_name 'my-sam-app', got %s", cfg.SAM.TransformOptions.StackName)
	}
}

func TestLoad_SAMConfigDefaults(t *testing.T) {
	content := `
templates:
  - "*.yaml"
sam:
  auto_transform: false
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, writeErr := tmpFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write temp file: %v", writeErr)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.SAM == nil {
		t.Fatal("Expected SAM config to be non-nil")
	}

	if cfg.SAM.AutoTransform {
		t.Error("Expected AutoTransform to be false")
	}
}

func TestMerge_SAMConfig(t *testing.T) {
	base := &Config{
		Templates: []string{"base/*.yaml"},
		SAM: &SAMConfig{
			AutoTransform: false,
			TransformOptions: &SAMTransformOptions{
				Region: "us-east-1",
			},
		},
	}

	override := &Config{
		Templates: []string{"override/*.yaml"},
		SAM: &SAMConfig{
			AutoTransform: true,
			TransformOptions: &SAMTransformOptions{
				Region:    "eu-west-1",
				AccountID: "123456789012",
			},
		},
	}

	result := Merge(base, override)

	if result.SAM == nil {
		t.Fatal("Expected SAM config in result")
	}

	if !result.SAM.AutoTransform {
		t.Error("Expected AutoTransform to be true from override")
	}

	if result.SAM.TransformOptions == nil {
		t.Fatal("Expected TransformOptions in result")
	}

	if result.SAM.TransformOptions.Region != "eu-west-1" {
		t.Errorf("Expected region 'eu-west-1', got %s", result.SAM.TransformOptions.Region)
	}

	if result.SAM.TransformOptions.AccountID != "123456789012" {
		t.Errorf("Expected account_id '123456789012', got %s", result.SAM.TransformOptions.AccountID)
	}
}

func TestMerge_SAMConfigNilBase(t *testing.T) {
	base := &Config{
		Templates: []string{"base/*.yaml"},
	}

	override := &Config{
		SAM: &SAMConfig{
			AutoTransform: true,
		},
	}

	result := Merge(base, override)

	if result.SAM == nil {
		t.Fatal("Expected SAM config in result")
	}

	if !result.SAM.AutoTransform {
		t.Error("Expected AutoTransform to be true from override")
	}
}

func TestMerge_SAMConfigNilOverride(t *testing.T) {
	base := &Config{
		Templates: []string{"base/*.yaml"},
		SAM: &SAMConfig{
			AutoTransform: true,
			TransformOptions: &SAMTransformOptions{
				Region: "us-east-1",
			},
		},
	}

	override := &Config{
		Templates: []string{"override/*.yaml"},
	}

	result := Merge(base, override)

	if result.SAM == nil {
		t.Fatal("Expected SAM config in result from base")
	}

	if !result.SAM.AutoTransform {
		t.Error("Expected AutoTransform to be true from base")
	}

	if result.SAM.TransformOptions.Region != "us-east-1" {
		t.Errorf("Expected region 'us-east-1', got %s", result.SAM.TransformOptions.Region)
	}
}
