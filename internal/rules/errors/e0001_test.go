package errors

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE0001_NoErrors(t *testing.T) {
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

	rule := &E0001{}
	matches := rule.Match(tmpl)

	// E0001 is a placeholder for transform errors caught at parse time
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestE0001_Metadata(t *testing.T) {
	rule := &E0001{}

	if rule.ID() != "E0001" {
		t.Errorf("Expected ID E0001, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
