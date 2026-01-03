package informational

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI1003_ShortDescription(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Description: Short description
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1003{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for short description, got %d", len(matches))
	}
}

func TestI1003_LongDescription(t *testing.T) {
	// Create a description that exceeds the warning threshold
	longDesc := strings.Repeat("This is a very long description. ", 30)
	yaml := "AWSTemplateFormatVersion: '2010-09-09'\nDescription: " + longDesc + "\nResources:\n  MyBucket:\n    Type: AWS::S3::Bucket\n"

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1003{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for long description approaching limit")
	}
}

func TestI1003_NoDescription(t *testing.T) {
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

	rule := &I1003{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for template without description, got %d", len(matches))
	}
}

func TestI1003_Metadata(t *testing.T) {
	rule := &I1003{}

	if rule.ID() != "I1003" {
		t.Errorf("Expected ID I1003, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
