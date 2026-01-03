package informational

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI7002_ShortMappingName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I7002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for short mapping name, got %d", len(matches))
	}
}

func TestI7002_LongMappingName(t *testing.T) {
	// Create a mapping name that exceeds the warning threshold (204+ characters)
	longName := "Mapping" + strings.Repeat("VeryLongName", 20)
	yaml := "AWSTemplateFormatVersion: '2010-09-09'\nMappings:\n  " + longName + ":\n    us-east-1:\n      AMI: ami-12345\nResources:\n  MyBucket:\n    Type: AWS::S3::Bucket\n"

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I7002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for long mapping name approaching limit")
	}
}

func TestI7002_Metadata(t *testing.T) {
	rule := &I7002{}

	if rule.ID() != "I7002" {
		t.Errorf("Expected ID I7002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
