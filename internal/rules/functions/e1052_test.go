package functions

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1052_ValidSSMRefs(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PARAM1: "{{resolve:ssm:/my/parameter}}"
          PARAM2: "{{resolve:ssm:/my/parameter:1}}"
          PARAM3: "{{resolve:ssm-secure:/my/secure/param}}"
          PARAM4: "{{resolve:ssm-secure:/my/secure/param:5}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid SSM references, got %d: %v", len(matches), matches)
	}
}

func TestE1052_MissingParameterName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PARAM: "{{resolve:ssm:}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing parameter name, got %d", len(matches))
	}
}

func TestE1052_InvalidParameterName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PARAM: "{{resolve:ssm:invalid param name}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid parameter name with spaces, got %d", len(matches))
	}
}

func TestE1052_InvalidVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PARAM: "{{resolve:ssm:/my/param:latest}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for non-integer version, got %d", len(matches))
	}
	if len(matches) > 0 && !strings.Contains(matches[0].Message, "must be an integer") {
		t.Errorf("Expected message about integer version, got: %s", matches[0].Message)
	}
}

func TestE1052_TrailingBackslash(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PARAM: "{{resolve:ssm:/my/param\\}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for trailing backslash, got %d", len(matches))
	}
}

func TestE1052_SSMInOutput(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  ParamValue:
    Value: "{{resolve:ssm:/my/parameter}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1052{}
	matches := rule.Match(parsed)

	// SSM refs ARE allowed in outputs (unlike secretsmanager)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for SSM ref in output (they are allowed), got %d: %v", len(matches), matches)
	}
}
