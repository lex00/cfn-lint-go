package errors

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2900_PlaceholderRule(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    Default: dev
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2900{}
	matches := rule.Match(tmpl)

	// E2900 is a placeholder rule that requires deployment file parsing infrastructure
	// It should return no matches when validating templates directly
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for placeholder rule E2900, got %d", len(matches))
	}
}

func TestE2900_Metadata(t *testing.T) {
	rule := &E2900{}

	if rule.ID() != "E2900" {
		t.Errorf("Expected ID E2900, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if rule.Source() == "" {
		t.Error("Source should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}

	// Verify it has the expected tags
	hasParametersTag := false
	hasDeploymentTag := false
	for _, tag := range tags {
		if tag == "parameters" {
			hasParametersTag = true
		}
		if tag == "deployment" {
			hasDeploymentTag = true
		}
	}

	if !hasParametersTag {
		t.Error("Expected 'parameters' tag")
	}

	if !hasDeploymentTag {
		t.Error("Expected 'deployment' tag")
	}
}
