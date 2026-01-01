package errors

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE0002_NoErrors(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E0002{}
	matches := rule.Match(tmpl)

	// E0002 is a placeholder for rule processing errors
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestE0002_Metadata(t *testing.T) {
	rule := &E0002{}

	if rule.ID() != "E0002" {
		t.Errorf("Expected ID E0002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
