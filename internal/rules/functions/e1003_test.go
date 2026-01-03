package functions

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1003_ValidDescription(t *testing.T) {
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

	rule := &E1003{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid description, got %d: %v", len(matches), matches)
	}
}

func TestE1003_DescriptionTooLong(t *testing.T) {
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

	rule := &E1003{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for description too long, got %d", len(matches))
	}
	if len(matches) > 0 && !strings.Contains(matches[0].Message, "1024 bytes") {
		t.Errorf("Expected message about 1024 bytes limit, got: %s", matches[0].Message)
	}
}

func TestE1003_DescriptionExactly1024(t *testing.T) {
	exactDesc := strings.Repeat("a", 1024)
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Description: "` + exactDesc + `"
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1003{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for exactly 1024 bytes, got %d", len(matches))
	}
}

func TestE1003_NoDescription(t *testing.T) {
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

	rule := &E1003{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no description, got %d: %v", len(matches), matches)
	}
}
