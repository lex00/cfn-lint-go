package parameters

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2010_UnderLimit(t *testing.T) {
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

	rule := &E2010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for under limit, got %d", len(matches))
	}
}

func TestE2010_OverLimit(t *testing.T) {
	// Generate template with 201 parameters
	var params []string
	for i := 0; i <= MaxParameters; i++ {
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

	rule := &E2010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for exceeding parameter limit")
	}
}

func TestE2010_ExactlyAtLimit(t *testing.T) {
	// Generate template with exactly 200 parameters
	var params []string
	for i := 0; i < MaxParameters; i++ {
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

	rule := &E2010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches at exactly limit, got %d", len(matches))
	}
}

func TestE2010_Metadata(t *testing.T) {
	rule := &E2010{}

	if rule.ID() != "E2010" {
		t.Errorf("Expected ID E2010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
