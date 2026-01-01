package functions

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1004_ValidDescription(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Description: "This is a valid description"
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1004{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid description, got %d: %v", len(matches), matches)
	}
}

func TestE1004_DescriptionTooLong(t *testing.T) {
	longDesc := strings.Repeat("a", 1025)
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Description: "` + longDesc + `"
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1004{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for description too long, got %d", len(matches))
	}
}

func TestE1004_NoDescription(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1004{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no description, got %d: %v", len(matches), matches)
	}
}
