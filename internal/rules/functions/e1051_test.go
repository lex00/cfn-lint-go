package functions

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1051_ValidSecretsManagerRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          DB_PASSWORD: "{{resolve:secretsmanager:my-secret:SecretString:password}}"
          API_KEY: "{{resolve:secretsmanager:api-key}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1051{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Secrets Manager references, got %d: %v", len(matches), matches)
	}
}

func TestE1051_MissingSecretID(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PASSWORD: "{{resolve:secretsmanager:}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1051{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing secret ID, got %d", len(matches))
	}
}

func TestE1051_InvalidInOutput(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  SecretValue:
    Value: "{{resolve:secretsmanager:my-secret}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1051{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for Secrets Manager ref in output, got %d", len(matches))
	}
	if len(matches) > 0 && !strings.Contains(matches[0].Message, "can only be used in resource properties") {
		t.Errorf("Expected message about resource properties only, got: %s", matches[0].Message)
	}
}

func TestE1051_TrailingBackslash(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          PASSWORD: "{{resolve:secretsmanager:my-secret\\}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1051{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for trailing backslash, got %d", len(matches))
	}
}
