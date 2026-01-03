package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1032_WithTransform(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Transform: AWS::LanguageExtensions
Resources:
  Fn::ForEach::Buckets:
    - BucketName
    - [Logs, Data]
    - ${BucketName}Bucket:
        Type: AWS::S3::Bucket
        Properties:
          BucketName: !Sub '${AWS::StackName}-${BucketName}'
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1032{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when AWS::LanguageExtensions transform is present, got %d: %v", len(matches), matches)
	}
}

func TestE1032_WithoutTransform(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  Fn::ForEach::Buckets:
    - BucketName
    - [Logs, Data]
    - ${BucketName}Bucket:
        Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1032{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match when AWS::LanguageExtensions transform is missing, got %d", len(matches))
	}
	if len(matches) > 0 && matches[0].Message == "" {
		t.Errorf("Expected non-empty error message")
	}
}

func TestE1032_NoFnForEach(t *testing.T) {
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

	rule := &E1032{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no Fn::ForEach is used, got %d", len(matches))
	}
}
