package informational

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI2011_ShortParameterName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  MyParam:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I2011{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for short parameter name, got %d", len(matches))
	}
}

func TestI2011_LongParameterName(t *testing.T) {
	// Create a parameter name that exceeds the warning threshold (204+ characters)
	longName := "Param" + strings.Repeat("VeryLongName", 20)
	yaml := "AWSTemplateFormatVersion: '2010-09-09'\nParameters:\n  " + longName + ":\n    Type: String\nResources:\n  MyBucket:\n    Type: AWS::S3::Bucket\n"

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I2011{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for long parameter name approaching limit")
	}
}

func TestI2011_Metadata(t *testing.T) {
	rule := &I2011{}

	if rule.ID() != "I2011" {
		t.Errorf("Expected ID I2011, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
