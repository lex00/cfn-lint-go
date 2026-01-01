package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8007_ValidConditionRef(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  IsProd:
    !Equals [!Ref Environment, prod]
  NotProd:
    !Not
      - Condition: IsProd

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

	rule := &E8007{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Condition intrinsic, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8007_ConditionNotString(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadCondition:
    Condition:
      - not-a-string

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8007{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Condition intrinsic with non-string value")
	}
}

func TestE8007_NoConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8007{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no conditions, got %d", len(matches))
	}
}

func TestE8007_Metadata(t *testing.T) {
	rule := &E8007{}

	if rule.ID() != "E8007" {
		t.Errorf("Expected ID E8007, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
