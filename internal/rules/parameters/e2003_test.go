package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2003_ValidParameterName(t *testing.T) {
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

	rule := &E2003{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid parameter name, got %d: %v", len(matches), matches)
	}
}

func TestE2003_InvalidParameterNameWithHyphen(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  my-parameter:
    Type: String
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2003{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid parameter name with hyphen, got %d", len(matches))
	}
}

func TestE2003_InvalidParameterNameStartsWithNumber(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  123Parameter:
    Type: String
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2003{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid parameter name starting with number, got %d", len(matches))
	}
}
