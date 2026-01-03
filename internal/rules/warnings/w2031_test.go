package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW2031_ValidAllowedPatternWithAnchors(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  BucketName:
    Type: String
    AllowedPattern: "^[a-z0-9-]+$"
    Default: my-bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2031{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid AllowedPattern with anchors, got %d: %v", len(matches), matches)
	}
}

func TestW2031_DefaultNotMatchingPattern(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  BucketName:
    Type: String
    AllowedPattern: "[a-z0-9-]+"
    Default: My_Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2031{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for default not matching pattern, got 0")
	}
}
