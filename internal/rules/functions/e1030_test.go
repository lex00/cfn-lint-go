package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1030_WithTransform(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Transform: AWS::LanguageExtensions
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      Tags:
        - Key: Count
          Value: !Length
            - item1
            - item2
            - item3
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1030{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when AWS::LanguageExtensions transform is present, got %d: %v", len(matches), matches)
	}
}

func TestE1030_WithoutTransform(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      Tags:
        - Key: Count
          Value:
            Fn::Length:
              - item1
              - item2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1030{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match when AWS::LanguageExtensions transform is missing, got %d", len(matches))
	}
	if len(matches) > 0 && matches[0].Message == "" {
		t.Errorf("Expected non-empty error message")
	}
}

func TestE1030_WithMultipleTransforms(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Transform:
  - AWS::Serverless-2016-10-31
  - AWS::LanguageExtensions
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Environment:
        Variables:
          COUNT:
            Fn::Length:
              - a
              - b
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1030{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when AWS::LanguageExtensions is in transform array, got %d: %v", len(matches), matches)
	}
}

func TestE1030_NoFnLength(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1030{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no Fn::Length is used, got %d", len(matches))
	}
}
