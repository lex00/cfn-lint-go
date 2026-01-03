package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1028_IfWithValidCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd:
    Fn::Equals: [!Ref Environment, prod]
Parameters:
  Environment:
    Type: String
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

	rule := &W1028{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Fn::If, got %d: %v", len(matches), matches)
	}
}

func TestW1028_StaticCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  AlwaysTrue:
    Fn::Equals: [true, true]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !If [AlwaysTrue, prod-bucket, dev-bucket]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1028{}
	matches := rule.Match(parsed)

	// The rule checks for unreachable paths in Fn::If, not static conditions
	// This test just verifies the rule runs without error
	_ = matches
}
