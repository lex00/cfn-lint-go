package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1040_ValidGetAttString(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement: []
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !GetAtt MyRole.Arn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid GetAtt string format, got %d: %v", len(matches), matches)
	}
}

func TestE1040_ValidGetAttArray(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement: []
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt:
          - MyRole
          - Arn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid GetAtt array format, got %d: %v", len(matches), matches)
	}
}

func TestE1040_MissingDot(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt: MyRoleArn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing dot in GetAtt string, got %d", len(matches))
	}
}

func TestE1040_EmptyResourceName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt: ".Arn"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for empty resource name, got %d", len(matches))
	}
}

func TestE1040_EmptyAttributeName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt: "MyRole."
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for empty attribute name, got %d", len(matches))
	}
}

func TestE1040_ArrayTooShort(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt:
          - MyRole
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1040{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for array too short, got %d", len(matches))
	}
}
