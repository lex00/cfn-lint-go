package informational

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI6011_ShortOutputName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I6011{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for short output name, got %d", len(matches))
	}
}

func TestI6011_LongOutputName(t *testing.T) {
	// Create an output name that exceeds the warning threshold (204+ characters)
	longName := "Output" + strings.Repeat("VeryLongName", 20)
	yaml := "AWSTemplateFormatVersion: '2010-09-09'\nResources:\n  MyBucket:\n    Type: AWS::S3::Bucket\nOutputs:\n  " + longName + ":\n    Value: !Ref MyBucket\n"

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I6011{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for long output name approaching limit")
	}
}

func TestI6011_Metadata(t *testing.T) {
	rule := &I6011{}

	if rule.ID() != "I6011" {
		t.Errorf("Expected ID I6011, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
