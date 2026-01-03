package informational

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3012_ShortResourceName(t *testing.T) {
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

	rule := &I3012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for short resource name, got %d", len(matches))
	}
}

func TestI3012_LongResourceName(t *testing.T) {
	// Create a resource name that exceeds the warning threshold (204+ characters)
	longName := "Resource" + strings.Repeat("VeryLongName", 20)
	yaml := "AWSTemplateFormatVersion: '2010-09-09'\nResources:\n  " + longName + ":\n    Type: AWS::S3::Bucket\n"

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for long resource name approaching limit")
	}
}

func TestI3012_Metadata(t *testing.T) {
	rule := &I3012{}

	if rule.ID() != "I3012" {
		t.Errorf("Expected ID I3012, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
