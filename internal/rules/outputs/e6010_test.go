package outputs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE6010_UnderLimit(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
  Output1:
    Value: value1
  Output2:
    Value: value2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for under limit, got %d", len(matches))
	}
}

func TestE6010_OverLimit(t *testing.T) {
	// Generate template with 201 outputs
	var outputs []string
	for i := 0; i <= MaxOutputs; i++ {
		outputs = append(outputs, fmt.Sprintf("  Output%d:\n    Value: value%d", i, i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
%s
`, strings.Join(outputs, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for exceeding output limit")
	}
}

func TestE6010_ExactlyAtLimit(t *testing.T) {
	// Generate template with exactly 200 outputs
	var outputs []string
	for i := 0; i < MaxOutputs; i++ {
		outputs = append(outputs, fmt.Sprintf("  Output%d:\n    Value: value%d", i, i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
%s
`, strings.Join(outputs, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches at exactly limit, got %d", len(matches))
	}
}

func TestE6010_Metadata(t *testing.T) {
	rule := &E6010{}

	if rule.ID() != "E6010" {
		t.Errorf("Expected ID E6010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
