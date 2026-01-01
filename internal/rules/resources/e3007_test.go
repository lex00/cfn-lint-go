package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3007_UniqueNames(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
  Bucket2:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3007{}
	matches := rule.Match(tmpl)

	// YAML parser enforces unique keys
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unique names, got %d", len(matches))
	}
}

func TestE3007_Metadata(t *testing.T) {
	rule := &E3007{}

	if rule.ID() != "E3007" {
		t.Errorf("Expected ID E3007, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
