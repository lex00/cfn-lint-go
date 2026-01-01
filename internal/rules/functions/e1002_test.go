package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1002_ValidSize(t *testing.T) {
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

	rule := &E1002{}
	matches := rule.Match(tmpl)

	// Size checking is done at parse time, not in Match
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestE1002_Metadata(t *testing.T) {
	rule := &E1002{}

	if rule.ID() != "E1002" {
		t.Errorf("Expected ID E1002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE1002_Constants(t *testing.T) {
	if TemplateSizeLimitDirect != 51200 {
		t.Errorf("Expected TemplateSizeLimitDirect to be 51200, got %d", TemplateSizeLimitDirect)
	}
	if TemplateSizeLimitS3 != 460800 {
		t.Errorf("Expected TemplateSizeLimitS3 to be 460800, got %d", TemplateSizeLimitS3)
	}
}
