package metadata

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE4001_ValidInterface(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  MyParam:
    Type: String
Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      - Label:
          default: General
        Parameters:
          - MyParam
    ParameterLabels:
      MyParam:
        default: My Parameter
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E4001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid interface, got %d: %v", len(matches), matches)
	}
}

func TestE4001_UndefinedParameterInGroup(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  MyParam:
    Type: String
Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      - Label:
          default: General
        Parameters:
          - NonExistentParam
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E4001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined parameter in group, got %d", len(matches))
	}
}

func TestE4001_UndefinedParameterInLabels(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  MyParam:
    Type: String
Metadata:
  AWS::CloudFormation::Interface:
    ParameterLabels:
      NonExistentParam:
        default: Invalid Parameter
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E4001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined parameter in labels, got %d", len(matches))
	}
}

func TestE4001_NoInterface(t *testing.T) {
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

	rule := &E4001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no interface, got %d: %v", len(matches), matches)
	}
}
