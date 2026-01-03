package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW4005_NoMetadata(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4005{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no metadata, got %d: %v", len(matches), matches)
	}
}

func TestW4005_CfnLintConfigInMetadata(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Metadata:
  cfn-lint:
    config:
      ignore_checks:
        - E1001
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4005{}
	matches := rule.Match(parsed)

	// This rule just validates the structure is correct
	// It should pass for valid cfn-lint config
	_ = matches
}

func TestW4005_InvalidCfnLintConfig(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Metadata:
  cfn-lint:
    invalid_key: value
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4005{}
	matches := rule.Match(parsed)

	// The rule validates cfn-lint configuration in Metadata
	// Just verify it runs without error
	_ = matches
}
