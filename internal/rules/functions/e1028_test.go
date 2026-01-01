package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1028_ValidFnIf(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Conditions:
  IsProd:
    Fn::Equals:
      - !Ref Environment
      - prod
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !If [IsProd, prod-bucket, dev-bucket]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1028{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Fn::If, got %d: %v", len(matches), matches)
	}
}

func TestE1028_UndefinedCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !If [NonExistentCondition, prod-bucket, dev-bucket]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1028{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined condition, got %d", len(matches))
	}
}

func TestE1028_WrongElementCount(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd:
    Fn::Equals:
      - true
      - true
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName:
        Fn::If:
          - IsProd
          - prod-bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1028{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for wrong element count, got %d", len(matches))
	}
}

func TestE1028_NestedFnIf(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd:
    Fn::Equals:
      - true
      - true
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !If
        - IsProd
        - !If [NonExistent, a, b]
        - dev-bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1028{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for nested undefined condition, got %d", len(matches))
	}
}

func TestE1028_InOutput(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !If [NonExistent, prod, dev]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1028{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined condition in output, got %d", len(matches))
	}
}
