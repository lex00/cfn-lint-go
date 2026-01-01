package parameters

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2011_ValidParameterNameLength(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  MyParameter:
    Type: String
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2011{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid parameter name length, got %d: %v", len(matches), matches)
	}
}

func TestE2011_ParameterNameTooLong(t *testing.T) {
	longName := strings.Repeat("A", 256)
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  ` + longName + `:
    Type: String
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2011{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for parameter name too long, got %d", len(matches))
	}
}
