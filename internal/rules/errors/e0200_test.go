package errors

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE0200_Metadata(t *testing.T) {
	rule := &E0200{}

	if rule.ID() != "E0200" {
		t.Errorf("Expected ID 'E0200', got '%s'", rule.ID())
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
}

func TestE0200_NoMatchesOnTemplate(t *testing.T) {
	// E0200 is for parameter file validation, not template validation
	// It should return no matches when checking a CloudFormation template
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  MyParam:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E0200{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for template (E0200 validates parameter files, not templates), got %d", len(matches))
	}
}
