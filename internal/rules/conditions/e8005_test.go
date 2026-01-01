package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8005_ValidNot(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  NotProd:
    !Not
      - !Equals [!Ref Environment, prod]

Parameters:
  Environment:
    Type: String

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8005{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Fn::Not, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8005_TooManyConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadNot:
    Fn::Not:
      - !Equals [a, a]
      - !Equals [b, b]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8005{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::Not with too many conditions")
	}
}

func TestE8005_EmptyList(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadNot:
    Fn::Not: []

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8005{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::Not with empty list")
	}
}

func TestE8005_Metadata(t *testing.T) {
	rule := &E8005{}

	if rule.ID() != "E8005" {
		t.Errorf("Expected ID E8005, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
