package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1050_ValidDynamicRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MySecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: my-secret
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          DB_PASSWORD: "{{resolve:secretsmanager:my-secret:SecretString:password}}"
          PARAM_VALUE: "{{resolve:ssm:/my/parameter}}"
          SECURE_PARAM: "{{resolve:ssm-secure:/my/secure/parameter}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid dynamic references, got %d: %v", len(matches), matches)
	}
}

func TestE1050_InvalidService(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          VALUE: "{{resolve:invalidservice:something}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid service, got %d", len(matches))
	}
}

func TestE1050_MissingParameterName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          VALUE: "{{resolve:ssm:}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing parameter name, got %d", len(matches))
	}
}

func TestE1050_MissingSecretId(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          VALUE: "{{resolve:secretsmanager:}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing secret ID, got %d", len(matches))
	}
}

func TestE1050_InOutput(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  SecretValue:
    Value: "{{resolve:invalid:ref}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid dynamic ref in output, got %d", len(matches))
	}
}

func TestE1050_SSMWithVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          VALUE: "{{resolve:ssm:/my/parameter:3}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1050{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for SSM with version, got %d: %v", len(matches), matches)
	}
}
