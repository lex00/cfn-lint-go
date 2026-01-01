package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1020_ValidRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref Environment
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1020{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Ref, got %d: %v", len(matches), matches)
	}
}

func TestE1020_RefIsArray(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName:
        Ref:
          - Environment
          - Something
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1020{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for Ref with array value, got %d", len(matches))
	}
}

func TestE1020_RefIsObject(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName:
        Ref:
          key: value
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1020{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for Ref with object value, got %d", len(matches))
	}
}

func TestE1020_RefInCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Conditions:
  IsProd:
    Fn::Equals:
      - Ref:
          - Environment
      - prod
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1020{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Ref in condition, got %d", len(matches))
	}
}
