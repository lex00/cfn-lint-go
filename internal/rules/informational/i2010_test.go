package informational

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI2010_LowParameterCount(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Param1:
    Type: String
  Param2:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I2010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for low parameter count, got %d", len(matches))
	}
}

func TestI2010_ApproachingLimit(t *testing.T) {
	// Generate template with 165 parameters (82.5% of limit)
	var params []string
	for i := 0; i < 165; i++ {
		params = append(params, fmt.Sprintf("  Param%d:\n    Type: String", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
%s
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`, strings.Join(params, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I2010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for parameter count approaching limit")
	}
}

func TestI2010_JustBelowThreshold(t *testing.T) {
	// Generate template with 159 parameters (just below 80% threshold)
	var params []string
	for i := 0; i < 159; i++ {
		params = append(params, fmt.Sprintf("  Param%d:\n    Type: String", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
%s
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`, strings.Join(params, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I2010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches just below threshold, got %d", len(matches))
	}
}

func TestI2010_Metadata(t *testing.T) {
	rule := &I2010{}

	if rule.ID() != "I2010" {
		t.Errorf("Expected ID I2010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
