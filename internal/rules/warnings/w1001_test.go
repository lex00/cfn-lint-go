package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1001_NoConditionalResources(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-conditional resources, got %d: %v", len(matches), matches)
	}
}

func TestW1001_RefToConditionalResource(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  CreateBucket:
    Fn::Equals: [true, true]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: CreateBucket
  MyPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for Ref to conditional resource, got %d", len(matches))
	}
}

func TestW1001_SameCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  CreateBucket:
    Fn::Equals: [true, true]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: CreateBucket
  MyPolicy:
    Type: AWS::S3::BucketPolicy
    Condition: CreateBucket
    Properties:
      Bucket: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when both resources have same condition, got %d: %v", len(matches), matches)
	}
}
