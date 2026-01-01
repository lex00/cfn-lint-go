package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW8001_UsedCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  CreateBucket:
    Type: String
Conditions:
  ShouldCreate:
    Fn::Equals: [!Ref CreateBucket, "true"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: ShouldCreate
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for used condition, got %d: %v", len(matches), matches)
	}
}

func TestW8001_UnusedCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  CreateBucket:
    Type: String
Conditions:
  ShouldCreate:
    Fn::Equals: [!Ref CreateBucket, "true"]
  UnusedCondition:
    Fn::Equals: ["a", "b"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: ShouldCreate
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for unused condition, got %d", len(matches))
	}
}

func TestW8001_ConditionUsedInFnIf(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  UseCustomName:
    Type: String
Conditions:
  CustomName:
    Fn::Equals: [!Ref UseCustomName, "true"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName:
        Fn::If: [CustomName, "custom-bucket", "default-bucket"]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for condition used in Fn::If, got %d: %v", len(matches), matches)
	}
}
